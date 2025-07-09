function Get-EthernetSpeeds {
  param (
    [string]$adapterName
  )
  # 監視間隔（秒）
  $interval = 1
  $samplingsNumber = 10

  # 各回の送信・受信速度（Kbps）を格納する配列
  $sentSpeeds = @()
  $receiveSpeeds = @()

  for ($i = 0; $i -lt $samplingsNumber; $i++) {
    # 開始時の統計情報取得
    $startStats = Get-NetAdapterStatistics -Name $adapterName
    Start-Sleep -Seconds $interval
    # 終了時の統計情報取得
    $endStats = Get-NetAdapterStatistics -Name $adapterName

    # 送信・受信バイト数の差分を計算
    $sentDiff = $endStats.SentBytes - $startStats.SentBytes
    $receivedDiff = $endStats.ReceivedBytes - $startStats.ReceivedBytes

    # 差分からKbpsへ変換 (1バイト=8ビット、1秒間隔の場合)
    $sentKbps = ($sentDiff * 8) / 1000
    $receivedKbps = ($receivedDiff * 8) / 1000

    # それぞれの速度を配列に追加
    $sentSpeeds += $sentKbps
    $receiveSpeeds += $receivedKbps
  }

  # 0でない値のみ抽出して平均を算出
  $nonZeroSent = $sentSpeeds | Where-Object { $_ -ne 0 }
  $nonZeroReceived = $receiveSpeeds | Where-Object { $_ -ne 0 }

  if ($nonZeroSent.Count -gt 0) {
    $avg_sent_kbps = ($nonZeroSent | Measure-Object -Average).Average
  } else {
    $avg_sent_kbps = 0
  }

  if ($nonZeroReceived.Count -gt 0) {
    $avg_received_kbps = ($nonZeroReceived | Measure-Object -Average).Average
  } else {
    $avg_received_kbps = 0
  }

  return @{
    "SentKbps" = [Math]::Round($avg_sent_kbps)
    "ReceivedKbps" = [Math]::Round($avg_received_kbps)
  }
}

function Get-SensorValue {
  param (
    [LibreHardwareMonitor.Hardware.Computer]$Monitor,
    [string]$SensorType,
    [string]$SensorName
  )
  $sensorValue = ""

  # Sensor一覧を取得
  $monitor.Hardware.Update()
  foreach ($hardware in $monitor.Hardware) {
    Write-Host "Hardware: {0} ---", $hardware.Name
    foreach ($subhardware in $hardware.SubHardware) {
      Write-Host "\tSubhardware: {0} ---", $subhardware.Name
      foreach ($sensor in $subhardware.Sensors) {
        Write-Host "\t\tsubhardwareSensor: ", $sensor.SensorType
        Write-Host "\t\tsubhardwareSensor: {0}, value: {1} ---", $sensor.Name, value: $sensor.Value
      }
    }
    foreach ($sensor in $hardware.Sensors) {
      Write-Host "\thardwareSensor: ", $sensor.SensorType
      Write-Host "\thardwareSensor: {0}, value: {1} ---", $sensor.Name, $sensor.Value
      if (($sensor.SensorType -eq $SensorType) -and ($sensor.Name -eq $SensorName)){
        $sensorValue = $sensor.Value
      }
    }
  }
  return $sensorValue
}

# CPUの温度取得
function Get-CPUTemperature {
  try {
    $dll = "$HOME\Downloads\picture_backup\lib\LibreHardwareMonitor-net472\LibreHardwareMonitorLib.dll"
    Add-Type -Path $dll

    $monitor = [LibreHardwareMonitor.Hardware.Computer]::new()
    $monitor.IsCPUEnabled = $true
    # $monitor.IsGPUEnabled = $true
    $monitor.Open()

    $cpuTemp = 0
    $interval = 1
    $samplingsNumber = 10
    $sensorValueTotal = 0

    for ($i = 0; $i -lt $samplingsNumber; $i++) {
      $sensorValue = Get-SensorValue $monitor "Temperature" "Core (Tctl/Tdie)"
      $sensorValueTotal += $sensorValue
      Start-Sleep -Seconds $interval
    }
    $monitor.Close()
    $sensorValueAvg = [math]::Round($sensorValueTotal / $samplingsNumber)
    return $sensorValueAvg
  } catch {
    Write-Host "エラー: $_"
    return $null
  }
}

# CPUの現在のクロック周波数取得
function Get-CPUClockSpeed {
  try {
    # WMIを使用してCPUクロック周波数を取得
    $cpuInfo = Get-WmiObject -Class Win32_Processor -ErrorAction Stop

    if ($cpuInfo) {
      # 現在のクロック周波数（MHz）
      $currentClockSpeed = $cpuInfo.CurrentClockSpeed

      # Open Hardware Monitorで詳細情報を取得する（利用可能な場合）
      $ohm = Get-WmiObject -Namespace "root\OpenHardwareMonitor" -Class Sensor -ErrorAction SilentlyContinue
      if ($ohm) {
        $cpuClock = $ohm | Where-Object { $_.SensorType -eq "Clock" -and $_.Name -like "*CPU Core*" }
        if ($cpuClock) {
          return $cpuClock.Value
        }
      }
      return $currentClockSpeed
    }
    Write-Host "CPUクロック周波数を取得できませんでした"
    return $null
  } catch {
    Write-Host "エラー: $_"
    return $null
  }
}

# GPUの情報取得（NVIDIA GPU用）
function Get-NVIDIAGPUInfo {
  try {
    # NVIDIA-SMIコマンドが利用可能か確認
    $nvidiaSmiPath = "C:\Windows\System32\nvidia-smi.exe"
    if (-not (Test-Path $nvidiaSmiPath)) {
      $nvidiaSmiPath = "C:\Program Files\NVIDIA Corporation\NVSMI\nvidia-smi.exe"
    }

    # GPUの情報を初期化
    $gpuInfo = @{
      "Name" = $null
      "Temperature" = $null
      "ClockSpeed" = $null
      "MaxClockSpeed" = $null
      "ShaderCores" = $null
      "RTCores" = $null
      "TensorCores" = $null
    }

    # WMIからGPU名を取得（基本情報）
    $wmiGPU = Get-WmiObject -Query "SELECT * FROM Win32_VideoController WHERE Name LIKE '%NVIDIA%'" -ErrorAction SilentlyContinue
    if ($wmiGPU) {
      $gpuInfo.Name = $wmiGPU.Name
    }else{
      Write-Host "Name: WMIが利用できないため取得不可"
    }

    if (Test-Path $nvidiaSmiPath) {
      # GPUの詳細情報取得
      $gpuTemp = & $nvidiaSmiPath --query-gpu=temperature.gpu --format=csv,noheader
      $gpuClock = & $nvidiaSmiPath --query-gpu=clocks.gr --format=csv,noheader
      $gpuMaxClock = & $nvidiaSmiPath --query-gpu=clocks.max.gr --format=csv,noheader
      $gpuModel = & $nvidiaSmiPath --query-gpu=name --format=csv,noheader

      # GPU名を更新（NVIDIA SMIからの情報が優先）
      if ($gpuModel) {
        $gpuInfo.Name = $gpuModel.Trim()
      }else{
        Write-Host "Name: nvidia-smiが利用できないため取得不可"
      }

      # 温度と周波数情報を更新
      if ($gpuTemp) {
        $gpuTemp = $gpuTemp.Trim()
        $gpuInfo.Temperature = [int]$gpuTemp
      } else {
        Write-Host "Temperature: nvidia-smiが利用できないため取得不可"
      }
      if ($gpuClock) {
        $gpuClock = $gpuClock.Trim()
        $gpuClock = $gpuClock -replace " MHz", ""
        $gpuInfo.ClockSpeed = [int]$gpuClock
      } else {
        Write-Host "ClockSpeed: nvidia-smiが利用できないため取得不可"
      }
      if ($gpuMaxClock) {
        $gpuMaxClock = $gpuMaxClock.Trim()
        $gpuMaxClock = $gpuMaxClock -replace " MHz", ""
        $gpuInfo.MaxClockSpeed = [int]$gpuMaxClock
      } else {
        Write-Host "MaxClockSpeed: nvidia-smiが利用できないため取得不可"
      }

      # GPU情報の詳細な取得を試みる
      $gpuFullInfo = & $nvidiaSmiPath -q
      $isSuper = $gpuInfo.Name -match "Super"

      # Shaderコア数（CUDAコア数）の推定
      if ($gpuInfo.Name -match "(RTX|GTX)\s+(\d+)") {
        $model = $matches[2]
        # GPU世代に基づいてコア数を推定
        $coreEstimate = $null
        $rtCoreEstimate = $null
        $tensorCoreEstimate = $null

        # RTX 30シリーズ
        if ($gpuInfo.Name -match "RTX\s+30(\d{2})") {
          switch ($matches[1]) {
            "90" {
              $coreEstimate = 10496
              $rtCoreEstimate = 82
              $tensorCoreEstimate = 328
            }
            "80" {
              $coreEstimate = 8704
              $rtCoreEstimate = 68
              $tensorCoreEstimate = 272
            }
            "70" {
              $coreEstimate = 5888
              $rtCoreEstimate = 46
              $tensorCoreEstimate = 184
            }
            "60" {
              $coreEstimate = 3584
              $rtCoreEstimate = 28
              $tensorCoreEstimate = 112
            }
            "50" {
              $coreEstimate = 2560
              $rtCoreEstimate = 20
              $tensorCoreEstimate = 80
            }
          }
        }
        # RTX 40シリーズ
        elseif ($gpuInfo.Name -match "RTX\s+40(\d{2})") {
          switch ($matches[1]) {
            "90" {
              $coreEstimate = 16384
              $rtCoreEstimate = 128
              $tensorCoreEstimate = 512
            }
            "80" {
              if ($isSuper) {
                $coreEstimate = 10240
                $rtCoreEstimate = 80
                $tensorCoreEstimate = 320
              } else {
                $coreEstimate = 9728
                $rtCoreEstimate = 76
                $tensorCoreEstimate = 304
              }
            }
            "70" {
              if ($isSuper) {
                $coreEstimate = 7168
                $rtCoreEstimate = 56
                $tensorCoreEstimate = 224
              } else {
                $coreEstimate = 5888
                $rtCoreEstimate = 46
                $tensorCoreEstimate = 184
              }
            }
            "60" {
              if ($isSuper) {
                $coreEstimate = 3840
                $rtCoreEstimate = 30
                $tensorCoreEstimate = 120
              } else {
                $coreEstimate = 3072
                $rtCoreEstimate = 24
                $tensorCoreEstimate = 96
              }
          }
            "50" {
              $coreEstimate = 2560
              $rtCoreEstimate = 20
              $tensorCoreEstimate = 80
            }
          }
        }

        if ($coreEstimate) {
          Write-Host "ShaderCores: $coreEstimate (推定値)"
          $gpuInfo.ShaderCores = $coreEstimate
        } else {
          Write-Host "ShaderCores: 想定外のGPUモデルであるため取得不可"
        }

        if ($rtCoreEstimate) {
          Write-Host "RTCores: $rtCoreEstimate (推定値)"
          $gpuInfo.RTCores = $rtCoreEstimate
        } else {
          Write-Host "RTCores: 想定外のGPUモデルであるため取得不可"
        }

        if ($tensorCoreEstimate) {
          Write-Host "TensorCores: $tensorCoreEstimate (推定値)"
          $gpuInfo.TensorCores = $tensorCoreEstimate
        } else {
          Write-Host "TensorCores: 想定外のGPUモデルであるため取得不可"
        }
      }
      return $gpuInfo
    }
  } catch {
    Write-Host "エラー: $_"
    return @{
      "Temperature" = $null
      "ClockSpeed" = $null
      "MaxClockSpeed" = $null
    }
  }
}

# AMD GPUの情報取得
function Get-AMDGPUInfo {
  try {
    # AMDのGPU情報をWMIから取得
    $gpu = Get-WmiObject -Query "SELECT * FROM Win32_VideoController WHERE Name LIKE '%AMD%' OR Name LIKE '%Radeon%'" -ErrorAction Stop

    # GPUの情報を初期化
    $gpuInfo = @{
      "Name" = $null
      "Temperature" = $null
      "ClockSpeed" = $null
      "MaxClockSpeed" = $null
      "ShaderCores" = $null
      "RTCores" = $null
      "TensorCores" = $null
    }

    if ($gpu) {
      # GPU名を設定
      $gpuInfo.Name = $gpu.Name

      # Open Hardware Monitorが利用可能か確認
      $ohm = Get-WmiObject -Namespace "root\OpenHardwareMonitor" -Class Sensor -ErrorAction SilentlyContinue
      if ($ohm) {
        $gpuTemp = $ohm | Where-Object { $_.SensorType -eq "Temperature" -and $_.Name -like "*GPU*" }
        $gpuClock = $ohm | Where-Object { $_.SensorType -eq "Clock" -and $_.Name -like "*GPU Core*" }

        if ($gpuTemp) { $gpuInfo.Temperature = $gpuTemp.Value }
        if ($gpuClock) { $gpuInfo.ClockSpeed = $gpuClock.Value }
      }else{
        Write-Host "Temperature: Open Hardware Monitorが利用できないため取得不可"
        Write-Host "ClockSpeed: Open Hardware Monitorが利用できないため取得不可"
      }

      # 最大クロック周波数は直接取得できないため、推定値または固定値を使用
      $maxClock = $gpu.VideoProcessor -replace ".*Radeon\s+\w+\s+(\d+).*", "$1"
      if ($maxClock -match "^\d+$") {
        Write-Host "MaxClockSpeed: $maxClock MHz (推定値)"
        $gpuInfo.MaxClockSpeed = $maxClock
      }else{
        Write-Host "MaxClockSpeed: Open Hardware Monitorが利用できないため取得不可"
      }

      # コア数の推定 (一部のRadeonモデルに基づく推定値)
      if ($gpuInfo.Name -match "Radeon\s+(RX\s+)?(\d{4})") {
        $model = $matches[2]

        # Shader Core数（Compute Unit数）の推定
        $cuEstimate = $null
        # RX 6000シリーズ
        if ($model -match "6(\d{3})") {
          switch ($matches[1]) {
            "950" { $cuEstimate = 5120 }  # RX 6950 XT
            "900" { $cuEstimate = 5120 }  # RX 6900 XT
            "800" { $cuEstimate = 3840 }  # RX 6800 XT
            "750" { $cuEstimate = 3840 }  # RX 6750 XT
            "700" { $cuEstimate = 2560 }  # RX 6700 XT
            "650" { $cuEstimate = 2048 }  # RX 6650 XT
            "600" { $cuEstimate = 1792 }  # RX 6600 XT
            "500" { $cuEstimate = 1024 }  # RX 6500 XT
          }
        }
        # RX 7000シリーズ
        elseif ($model -match "7(\d{3})") {
          switch ($matches[1]) {
            "900" { $cuEstimate = 6144 }  # RX 7900 XTX
            "800" { $cuEstimate = 4608 }  # RX 7800 XT
            "700" { $cuEstimate = 3456 }  # RX 7700 XT
            "600" { $cuEstimate = 2048 }  # RX 7600 XT
          }
        }

        if ($cuEstimate) {
          Write-Host "ShaderCores: 推定値"
          $gpuInfo.ShaderCores = $cuEstimate
        } else {
          Write-Host "ShaderCores: Ray Acceleratorsが取得できなかったため、推定不可"
        }

        # RTコア数（Ray Accelerator数）の推定
        $raEstimate = $null
        if ($model -match "6(\d{3})") {
          switch ($matches[1]) {
            "950" { $raEstimate = 80 }  # RX 6950 XT
            "900" { $raEstimate = 80 }  # RX 6900 XT
            "800" { $raEstimate = 60 }  # RX 6800 XT
            "750" { $raEstimate = 60 }  # RX 6750 XT
            "700" { $raEstimate = 40 }  # RX 6700 XT
            "650" { $raEstimate = 32 }  # RX 6650 XT
            "600" { $raEstimate = 28 }  # RX 6600 XT
            "500" { $raEstimate = 16 }  # RX 6500 XT
          }
        }
        elseif ($model -match "7(\d{3})") {
          switch ($matches[1]) {
            "900" { $raEstimate = 96 }  # RX 7900 XTX
            "800" { $raEstimate = 72 }  # RX 7800 XT
            "700" { $raEstimate = 54 }  # RX 7700 XT
            "600" { $raEstimate = 32 }  # RX 7600 XT
          }
        }else{
          Write-Host "RTCores: 想定外のGPUモデルであるため取得不可"
        }

        if ($raEstimate) {
          Write-Host "RTCores: Ray Accelerators, 推定値"
          $gpuInfo.RTCores = $raEstimate
        } else {
          Write-Host "RTCores: Ray Acceleratorsが取得できなかったため、推定不可"
        }

        # TensorコアはAMD GPUには存在しない
        Write-Host "TensorCores: なし (AMDはAI Acceleratorsを使用)"
        $gpuInfo.TensorCores = $null
      }
    }

    return $gpuInfo
  } catch {
    Write-Host "エラー: $_"
    return @{
      "Name" = $null
      "Temperature" = $null
      "ClockSpeed" = $null
      "MaxClockSpeed" = $null
      "ShaderCores" = $null
      "RTCores" = $null
      "TensorCores" = $null
    }
  }
}

function Get-MachineInfo {
  # CPU情報の取得と変数への格納
  $cpuInfo = Get-CimInstance Win32_Processor
  $cpuName = $cpuInfo.Name.Trim()
  $cpuMaxClockSpeed = $cpuInfo.MaxClockSpeed
  $cpuMaxClockSpeedGHz = $cpuMaxClockSpeed / 1000
  $cpuCores = $cpuInfo.NumberOfCores
  $cpuLogicalProcessors = $cpuInfo.NumberOfLogicalProcessors

  # メモリ使用量の取得
  $osInfo = Get-CimInstance Win32_OperatingSystem
  $usedMemoryKB = $osInfo.TotalVisibleMemorySize - $osInfo.FreePhysicalMemory
  $usedMemoryMB = $usedMemoryKB / 1024
  $usedMemoryGB = $usedMemoryKB / 1024 / 1024

  # 標準メモリ情報の取得
  $memoryInfo = Get-CimInstance Win32_PhysicalMemory
  $memoryNames = $memoryInfo | ForEach-Object { $_.PartNumber.Trim() }
  $memoryNames = $memoryNames -join ','
  $memoryManufacturers = $memoryInfo | ForEach-Object { $_.Manufacturer.Trim() }
  $memoryManufacturers = $memoryManufacturers -join ','
  $totalMemoryMB = [Math]::Round(($memoryInfo | Measure-Object -Property Capacity -Sum).Sum / 1MB, 2)
  $totalMemoryGB = [Math]::Round(($memoryInfo | Measure-Object -Property Capacity -Sum).Sum / 1GB, 2)

  # 現在の実際のメモリ速度を取得（オーバークロック含む）
  try {
      # WMIでメモリの現在の動作周波数を取得
      $currentMemorySpeed = Get-CimInstance -Namespace "root\wmi" -ClassName "MSMemory_CurrentSpeed" -ErrorAction Stop | Select-Object -ExpandProperty CurrentSpeed
  } catch {
    # 上記の方法で取得できない場合はWindows Registry から取得を試みる
    try {
      $registryPath = "HKLM:\HARDWARE\DESCRIPTION\System\CentralProcessor\0"
      $currentMemorySpeed = (Get-ItemProperty -Path $registryPath -Name "~MHz" -ErrorAction Stop)."~MHz"
    } catch {
      # それでも取得できない場合は標準のメモリ速度を使用
      $currentMemorySpeed = $memoryInfo[0].Speed
      Write-Host "注意: 実際のメモリ速度を取得できませんでした。定格速度を表示します。"
    }
  }

  # 結果を表示
  Write-Host "Memory current usage (MB): " $usedMemoryMB.ToString("N2")
  Write-Host "CPU name: " $cpuName
  $log = "CPU Max Clock Speed: {0} MHz ({1} GHz)" -f $cpuMaxClockSpeed, $cpuMaxClockSpeedGHz.ToString("N2")
  Write-Host $log
  Write-Host "CPU Physical Cores: " $cpuCores
  Write-Host "CPU Logical Processors: " $cpuLogicalProcessors
  Write-Host "Memory name: " $memoryNames
  Write-Host "Memory manufacturers: " $memoryManufacturers
  Write-Host "Memory total size (GB): " $totalMemoryGB
  Write-Host "Memory current clock speed (MHz, overclocked): " $currentMemorySpeed

  # ホスト名の取得
  $hostname = $env:COMPUTERNAME
  Write-Host "ホスト名: $hostname"

  # 対象のアダプター名（例: "イーサネット"）
  $ethernetInfo = Get-EthernetSpeeds "イーサネット"

  # 結果の表示
  $log = "Average Sent Speed: {0} Kbps" -f $ethernetInfo.SentKbps
  Write-Host $log
  $log = "Average Received Speed: {0} Kbps" -f $ethernetInfo.ReceivedKbps
  Write-Host $log

  $cpuTemp = Get-CPUTemperature
  $cpuCurrentClockSpeed = Get-CPUClockSpeed
  $nvidiaInfo = Get-NVIDIAGPUInfo
  $amdInfo = Get-AMDGPUInfo

  # 取得した各値をオブジェクトにまとめる
  $data = [PSCustomObject]@{
    cpu_name = $cpuName
    cpu_current_clock_speed_mhz = $cpuCurrentClockSpeed
    cpu_max_clock_speed_mhz = $cpuMaxClockSpeed
    cpu_cores = $cpuCores
    cpu_logical_processors = $cpuLogicalProcessors
    cpu_temperature = $cpuTemp
    memory_names = $memoryNames
    memory_manufacturers = $memoryManufacturers
    memory_total_mb = $totalMemoryMB
    memory_current_speed_mhz = $currentMemorySpeed
    memory_usage_mb = [math]::Round($usedMemoryMB, 2)
    pc_hostname = $hostname
    ethernet_avg_sent_kbps = $ethernetInfo.SentKbps
    ethernet_avg_received_kbps = $ethernetInfo.ReceivedKbps
    gpu_temperature_nvidia = $nvidiaInfo.Temperature
    gpu_current_clock_speed_mhz_nvidia = $nvidiaInfo.ClockSpeed
    gpu_max_clock_speed_mhz_nvidia = $nvidiaInfo.MaxClockSpeed
    gpu_shader_cores_mhz_nvidia = $nvidiaInfo.ShaderCores
    gpu_rt_cores_mhz_nvidia = $nvidiaInfo.RTCores
    gpu_tensor_cores_mhz_nvidia = $nvidiaInfo.TensorCores
    gpu_temperature_amd = $amdInfo.Temperature
    gpu_current_clock_speed_mhz_amd = $amdInfo.ClockSpeed
    gpu_max_clock_speed_mhz_amd = $amdInfo.MaxClockSpeed
    gpu_shader_cores_mhz_amd = $amdInfo.ShaderCores
    gpu_rt_cores_mhz_amd = $amdInfo.RTCores
    gpu_tensor_cores_mhz_amd = $amdInfo.TensorCores
    output_at = [int][double]::Parse((Get-Date -UFormat %s))
  }

  return $data
}

$machineInfo = Get-MachineInfo

Write-Host "取得したシステム情報オブジェクト:"
Write-Host $machineInfo -ForegroundColor Green

$logJson = "log_{0}.json" -f (Get-Date).ToString("yyyyMMdd-HHmmss")

# オブジェクトを JSON 形式に変換 (-Depth パラメータでネストの深さを調整)
$json = $machineInfo | ConvertTo-Json -Depth 3
Write-Host "取得したシステム情報JSON:"
Write-Host $json -ForegroundColor Green
$json | Set-Content -Path "$HOME\Downloads\picture_backup\$logJson" -Encoding utf8

# # 送信先の API エンドポイント (実際の URL に置き換えてください)
# $apiUrl = "https://your-api-endpoint.example.com/post"

# # JSON データを API に POST する
# try {
#   $response = Invoke-RestMethod -Uri $apiUrl -Method Post -Body $json -ContentType "application/json"
#   Write-Host "データ送信成功。レスポンス:" $response
# } catch {
#   Write-Host "データ送信失敗:" $_.Exception.Message
# }

# # 結果の表示
# Write-Host "取得したシステム情報JSON:"
# Write-Host $json
