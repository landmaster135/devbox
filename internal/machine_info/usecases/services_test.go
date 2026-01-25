package usecases

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type mockFileSystem struct {
	ensureDirCalls []ensureDirCall
	writeFileCalls []writeFileCall
	joinCalls      [][]string
	ensureDirErr   error
	writeFileErr   error
}

type ensureDirCall struct {
	path string
	perm os.FileMode
}

type writeFileCall struct {
	path string
	data []byte
	perm os.FileMode
}

func (m *mockFileSystem) EnsureDir(path string, perm os.FileMode) error {
	m.ensureDirCalls = append(m.ensureDirCalls, ensureDirCall{path: path, perm: perm})
	return m.ensureDirErr
}

func (m *mockFileSystem) WriteFile(path string, data []byte, perm os.FileMode) error {
	m.writeFileCalls = append(m.writeFileCalls, writeFileCall{path: path, data: append([]byte(nil), data...), perm: perm})
	return m.writeFileErr
}

func (m *mockFileSystem) Join(elem ...string) string {
	m.joinCalls = append(m.joinCalls, append([]string{}, elem...))
	return filepath.Join(elem...)
}

func TestSaveMachineInfoLog_Success(t *testing.T) {
	mockFS := &mockFileSystem{}
	svc := NewMachineInfoServiceWithDependencies(mockFS)
	info := &MachineInfo{
		CPUName:                 "Test CPU",
		CPUCurrentClockSpeedMHz: 3000.25,
		CPUMaxClockSpeedMHz:     4200.75,
		CPUCores:                8,
		CPULogicalProcessors:    16,
		CPUTemperature:          55.5,
		MemoryTotalMB:           32000,
		MemoryUsageMB:           10240,
		PCHostname:              "devbox",
		EthernetAvgSentKbps:     1000.5,
		EthernetAvgReceivedKbps: 2000.25,
		OutputAt:                1234567890,
	}

	jsonText, outputPath, err := svc.SaveMachineInfoLog(info, " /tmp/machine-info ")
	if err != nil {
		t.Fatalf("SaveMachineInfoLog returned error: %v", err)
	}

	if len(mockFS.ensureDirCalls) != 1 {
		t.Fatalf("expected EnsureDir to be called once, got %d", len(mockFS.ensureDirCalls))
	}
	if got := mockFS.ensureDirCalls[0].path; got != "/tmp/machine-info" {
		t.Fatalf("expected EnsureDir path '/tmp/machine-info', got %s", got)
	}
	if mockFS.ensureDirCalls[0].perm != 0o755 {
		t.Fatalf("expected EnsureDir perm 0755, got %v", mockFS.ensureDirCalls[0].perm)
	}

	if len(mockFS.joinCalls) != 1 {
		t.Fatalf("expected Join to be called once, got %d", len(mockFS.joinCalls))
	}
	joinArgs := mockFS.joinCalls[0]
	if len(joinArgs) != 2 {
		t.Fatalf("expected Join to receive two arguments, got %d", len(joinArgs))
	}
	if joinArgs[0] != "/tmp/machine-info" {
		t.Fatalf("expected first Join arg '/tmp/machine-info', got %s", joinArgs[0])
	}
	if !strings.HasPrefix(joinArgs[1], "log_") || !strings.HasSuffix(joinArgs[1], ".json") {
		t.Fatalf("unexpected filename generated: %s", joinArgs[1])
	}

	if len(mockFS.writeFileCalls) != 1 {
		t.Fatalf("expected WriteFile to be called once, got %d", len(mockFS.writeFileCalls))
	}
	writeCall := mockFS.writeFileCalls[0]
	if writeCall.perm != 0o644 {
		t.Fatalf("expected WriteFile perm 0644, got %v", writeCall.perm)
	}
	if writeCall.path != outputPath {
		t.Fatalf("expected output path %s, got %s", outputPath, writeCall.path)
	}
	if string(writeCall.data) != jsonText {
		t.Fatalf("expected written JSON to match returned text")
	}

	var decoded MachineInfo
	if err := json.Unmarshal([]byte(jsonText), &decoded); err != nil {
		t.Fatalf("failed to decode JSON: %v", err)
	}
	if decoded.CPUName != info.CPUName || decoded.OutputAt != info.OutputAt {
		t.Fatalf("decoded JSON does not match original info: %#v", decoded)
	}
}

func TestSaveMachineInfoLog_Errors(t *testing.T) {
	svc := NewMachineInfoServiceWithDependencies(&mockFileSystem{})
	if _, _, err := svc.SaveMachineInfoLog(nil, "."); err == nil {
		t.Fatalf("expected error when info is nil")
	}

	mEnsureErr := errors.New("ensure failure")
	svc = NewMachineInfoServiceWithDependencies(&mockFileSystem{ensureDirErr: mEnsureErr})
	_, _, err := svc.SaveMachineInfoLog(&MachineInfo{CPUName: "cpu"}, "/tmp")
	if !errors.Is(err, mEnsureErr) {
		t.Fatalf("expected ensure dir error, got %v", err)
	}

	mWriteErr := errors.New("write failure")
	svc = NewMachineInfoServiceWithDependencies(&mockFileSystem{writeFileErr: mWriteErr})
	_, _, err = svc.SaveMachineInfoLog(&MachineInfo{CPUName: "cpu"}, "/tmp")
	if !errors.Is(err, mWriteErr) {
		t.Fatalf("expected write file error, got %v", err)
	}
}

type mockCollectorSaver struct {
	ifaceCalls  []string
	outputCalls []string
	savedInfos  []*MachineInfo
	result      *MachineInfoResult
	collectErr  error
	saveJSON    string
	savePath    string
	saveErr     error
}

func (m *mockCollectorSaver) CollectUbuntuInfo(networkInterface string) (*MachineInfoResult, error) {
	m.ifaceCalls = append(m.ifaceCalls, networkInterface)
	if m.collectErr != nil {
		return nil, m.collectErr
	}
	return m.result, nil
}

func (m *mockCollectorSaver) SaveMachineInfoLog(info *MachineInfo, outputDir string) (string, string, error) {
	m.savedInfos = append(m.savedInfos, info)
	m.outputCalls = append(m.outputCalls, outputDir)
	if m.saveErr != nil {
		return "", "", m.saveErr
	}
	return m.saveJSON, m.savePath, nil
}

func (m *mockCollectorSaver) CollectAndSaveUbuntuInfo(networkInterface, outputDir string) (*MachineInfoResult, string, string, error) {
	return collectAndSave(m, networkInterface, outputDir)
}

func TestCollectAndSaveHelperSuccess(t *testing.T) {
	mock := &mockCollectorSaver{
		result:   &MachineInfoResult{Info: &MachineInfo{CPUName: "Mock"}},
		saveJSON: `{"cpu_name":"Mock"}`,
		savePath: "/tmp/log.json",
	}

	result, jsonText, outputPath, err := collectAndSave(mock, "eth1", "/tmp")
	if err != nil {
		t.Fatalf("collectAndSave returned error: %v", err)
	}
	if len(mock.ifaceCalls) != 1 || mock.ifaceCalls[0] != "eth1" {
		t.Fatalf("expected collector to be called with eth1, got %v", mock.ifaceCalls)
	}
	if len(mock.outputCalls) != 1 || mock.outputCalls[0] != "/tmp" {
		t.Fatalf("expected saver to be called with /tmp, got %v", mock.outputCalls)
	}
	if len(mock.savedInfos) != 1 || mock.savedInfos[0].CPUName != "Mock" {
		t.Fatalf("unexpected info saved: %#v", mock.savedInfos)
	}
	if result.Info.CPUName != "Mock" {
		t.Fatalf("unexpected result: %#v", result)
	}
	if jsonText != mock.saveJSON || outputPath != mock.savePath {
		t.Fatalf("unexpected outputs json=%s path=%s", jsonText, outputPath)
	}
}

func TestCollectAndSaveHelperErrors(t *testing.T) {
	mCollectErr := errors.New("collect fail")
	mock := &mockCollectorSaver{collectErr: mCollectErr}
	if _, _, _, err := collectAndSave(mock, "eth0", "/tmp"); !errors.Is(err, mCollectErr) {
		t.Fatalf("expected collect error, got %v", err)
	}

	mock = &mockCollectorSaver{result: &MachineInfoResult{}}
	if _, _, _, err := collectAndSave(mock, "eth0", "/tmp"); err == nil {
		t.Fatalf("expected error when info is nil")
	}

	saveErr := errors.New("save fail")
	mock = &mockCollectorSaver{
		result:   &MachineInfoResult{Info: &MachineInfo{CPUName: "Mock"}},
		saveErr:  saveErr,
		saveJSON: "",
	}
	if _, _, _, err := collectAndSave(mock, "eth0", "/tmp"); !errors.Is(err, saveErr) {
		t.Fatalf("expected save error, got %v", err)
	}
}

func TestParseProcMeminfoWithAvailable(t *testing.T) {
	meminfo := `MemTotal:        7909052 kB
MemAvailable:    2256340 kB
MemFree:          534360 kB
Buffers:          149000 kB
Cached:          1884336 kB
SReclaimable:     250000 kB
Shmem:            534360 kB
`

	total, used, err := parseProcMeminfo([]byte(meminfo))
	if err != nil {
		t.Fatalf("parseProcMeminfo returned error: %v", err)
	}

	if total != 7723 {
		t.Fatalf("unexpected total: got %d MB", total)
	}
	if used != 5520 {
		t.Fatalf("unexpected used: got %d MB", used)
	}
}

func TestParseProcMeminfoFallbackWithoutAvailable(t *testing.T) {
	meminfo := `MemTotal:        2048000 kB
MemFree:          512000 kB
Buffers:           64000 kB
Cached:           256000 kB
SReclaimable:      32000 kB
Shmem:             16000 kB
`

	total, used, err := parseProcMeminfo([]byte(meminfo))
	if err != nil {
		t.Fatalf("parseProcMeminfo returned error: %v", err)
	}

	if total != 2000 {
		t.Fatalf("unexpected total: got %d MB", total)
	}
	if used != 1171 {
		t.Fatalf("unexpected used: got %d MB", used)
	}
}
