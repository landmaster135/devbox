/**
 * script.js - 画像トリミング座標プロットツール
 * 画像の読み込みとマーカー操作の機能を実装
 */

// イベントリスナーの設定
const handleImageChange = (e) => {
  update(false, { isInputting: false, e: null }, { isUpdating: false, e: e }, { isImageChanged: true, e: e });
  updateLockStatusDisplay(true, false); // 画像読み込み後のロック状態を表示
};
document.getElementById("imageLoader").addEventListener("change", handleImageChange);

const handleBarsAndCoordinates = (e) => {
  const isLocked = updateLockStatusDisplay(false, null); // マウス移動時にロック状態を更新
  if (!isLocked) { // ロックされていない場合のみバーを移動
    update(false, { isInputting: false, e: null }, { isUpdating: true, e: e }, { isImageChanged: false, e: null });
  }
};
document.addEventListener("mousemove", handleBarsAndCoordinates);

const updateCoordinateInput = (e) => {
  updateLockStatusDisplay(true, true); // ロック状態を更新
  update(false, { isInputting: true, e: e }, { isUpdating: false, e: e }, { isImageChanged: false, e: null });
}
document.addEventListener("mousedown", updateCoordinateInput);

const unlockBar = (e) => {
  updateLockStatusDisplay(true, false); // ロック状態を更新
  update(true, { isInputting: false, e: null }, { isUpdating: false, e: null }, { isImageChanged: false, e: null });
};
document.addEventListener("mouseup", unlockBar);

// ロック状態の表示を更新する関数
const updateLockStatusDisplay = (isUpdating, isLocking) => {
  let flg;
  if (isUpdating) {
    flg = isLocking;
  }

  const lockStatusDisplay = document.getElementById("lockStatus");
  lockStatusDisplay.textContent = flg ? "ロック中..." : "解除中...";
  lockStatusDisplay.style.color = flg ? "#e74c3c" : "#2ecc71"; // 赤/緑で状態を表示
  return flg;
};

// クリップボードに座標をコピーする関数
const copyToClipboard = (copyStatusesToCb, clipboardCoordinates, coordinateXUpper, coordinateYUpper, coordinateXLower, coordinateYLower) => {
  const sep = " ";
  // コマンドライン引数形式で座標を構築
  const copyText = `-x1 ${coordinateXUpper.value}${sep}-y1 ${coordinateYUpper.value}${sep}-x2 ${coordinateXLower.value}${sep}-y2 ${coordinateYLower.value}`;

  // クリップボード表示用テキストボックスに座標を表示
  clipboardCoordinates.value = copyText;

  // クリップボードにコピー
  navigator.clipboard.writeText(copyText).then(() => {
    // コピー成功時の表示
    copyStatusesToCb.textContent = "コピーしました！";
    setTimeout(() => {
      copyStatusesToCb.textContent = "";
    }, 2000);
  }, (err) => {
    // コピー失敗時の表示
    console.error("クリップボードへのコピーに失敗しました", err);
    copyStatusesToCb.textContent = "コピーに失敗しました";
    setTimeout(() => {
      copyStatusesToCb.textContent = "";
    }, 2000);
  });
};

// アクティブなバーのインデックスを管理する状態更新関数
const updateStateOfBarActivation = (indexOfActiveBar) => {
  const NUMBER_OF_BAR = 2; // バーの数（上部と下部）

  return (isBarStateUpdating, propsForCoordRecordInput, propsForUpdateCoordAndBar, propsForImageChanging) => {
    console.log("バーの状態更新を開始します");

    // アクティブなバーを更新
    if (isBarStateUpdating) {
      indexOfActiveBar += 1;
      if (indexOfActiveBar === NUMBER_OF_BAR) {
        indexOfActiveBar = 0;
      }
    }
    console.log(`アクティブなバー: ${indexOfActiveBar}`);

    // DOM要素の取得
    const imageElement = document.getElementById("imageElement");
    const coordinateXUpperToRecord = document.getElementById("coordinateXUpperDisplayToRecord");
    const coordinateYUpperToRecord = document.getElementById("coordinateYUpperDisplayToRecord");
    const coordinateXLowerToRecord = document.getElementById("coordinateXLowerDisplayToRecord");
    const coordinateYLowerToRecord = document.getElementById("coordinateYLowerDisplayToRecord");
    const copyStatusDisplayToCb = document.getElementById("copyToCb");
    const clipboardCoordinates = document.getElementById("clipboardCoordinates");

    // 座標入力処理
    if (propsForCoordRecordInput.isInputting) {
      const imageRect = imageElement.getBoundingClientRect();
      const x = Math.round(propsForCoordRecordInput.e.clientX - imageRect.left);
      const y = Math.round(propsForCoordRecordInput.e.clientY - imageRect.top);

      // アクティブなバーに応じて座標を更新
      switch (indexOfActiveBar) {
        case 0: // 上部座標
          coordinateXUpperToRecord.value = `${x}`;
          coordinateYUpperToRecord.value = `${y}`;
          break;
        case 1: // 下部座標
          coordinateXLowerToRecord.value = `${x}`;
          coordinateYLowerToRecord.value = `${y}`;
          break;
      }

      // アクティブなバーに応じて更新対象を選択
      let coordXToRecord = coordinateXUpperToRecord;
      let coordYToRecord = coordinateYUpperToRecord;

      switch (indexOfActiveBar) {
        case 0: // 上部バー
          // デフォルト値を使用
          break;
        case 1: // 下部バー
          coordXToRecord = coordinateXLowerToRecord;
          coordYToRecord = coordinateYLowerToRecord;
          break;
      }

      const updateCoordinates = (imgRect, ev, coordX, coordY) => {
        // 画像内のマウス位置のみ処理
        let yPos = 0;
        if (ev.clientY >= imgRect.top && ev.clientY <= imgRect.bottom){
          yPos = ev.clientY - imgRect.top;
        } else if (ev.clientY > imgRect.bottom){
          yPos = imgRect.bottom - imgRect.top;
        }
        coordY.value = `${Math.round(yPos)}`;

        let xPos = 0;
        if (ev.clientX >= imgRect.left && ev.clientX <= imgRect.right) {
          xPos = ev.clientX - imgRect.left;
        } else if (ev.clientX > imgRect.right) {
          xPos = imgRect.right - imgRect.left;
        }
        coordX.value = `${Math.round(xPos)}`;
      };

      // バーと座標を更新
      updateCoordinates(
        imageRect,
        propsForUpdateCoordAndBar.e,
        coordXToRecord,
        coordYToRecord
      );

      // クリップボードに座標をコピー
      copyToClipboard(
        copyStatusDisplayToCb,
        clipboardCoordinates,
        coordinateXUpperToRecord,
        coordinateYUpperToRecord,
        coordinateXLowerToRecord,
        coordinateYLowerToRecord
      );
    }

    const coordinateXUpperDisplay = document.getElementById("coordinateXUpperDisplay");
    const coordinateYUpperDisplay = document.getElementById("coordinateYUpperDisplay");
    const coordinateXLowerDisplay = document.getElementById("coordinateXLowerDisplay");
    const coordinateYLowerDisplay = document.getElementById("coordinateYLowerDisplay");
    const draggableBarXUpper = document.getElementById("draggableBarXUpper");
    const draggableBarYUpper = document.getElementById("draggableBarYUpper");
    const draggableBarXLower = document.getElementById("draggableBarXLower");
    const draggableBarYLower = document.getElementById("draggableBarYLower");

    // バーと座標の更新処理
    if (propsForUpdateCoordAndBar.isUpdating) {
      const updateBarsAndCoordinates = (imgRect, ev, barX, barY, coordX, coordY) => {
        // 画像内のマウス位置のみ処理
        if (ev.clientX >= imgRect.left && ev.clientX <= imgRect.right &&
          ev.clientY >= imgRect.top && ev.clientY <= imgRect.bottom) {

          // Y座標の更新
          const yPos = ev.clientY - imgRect.top;
          barY.style.top = `${yPos}px`;
          barY.style.transform = 'translateY(-50%)'; // 中央揃え
          coordY.value = `${Math.round(yPos)}`;

          // X座標の更新
          const xPos = ev.clientX - imgRect.left;
          barX.style.left = `${xPos}px`;
          coordX.value = `${Math.round(xPos)}`;
        }
      };

      // 画像の位置情報を取得
      const imageRect = imageElement.getBoundingClientRect();

      // アクティブなバーに応じて更新対象を選択
      let draggableBarX = draggableBarXUpper;
      let draggableBarY = draggableBarYUpper;
      let coordXDisplay = coordinateXUpperDisplay;
      let coordYDisplay = coordinateYUpperDisplay;

      switch (indexOfActiveBar) {
        case 0: // 上部バー
          // デフォルト値を使用
          break;
        case 1: // 下部バー
          draggableBarX = draggableBarXLower;
          draggableBarY = draggableBarYLower;
          coordXDisplay = coordinateXLowerDisplay;
          coordYDisplay = coordinateYLowerDisplay;
          break;
      }

      // バーと座標を更新
      updateBarsAndCoordinates(
        imageRect,
        propsForUpdateCoordAndBar.e,
        draggableBarX,
        draggableBarY,
        coordXDisplay,
        coordYDisplay
      );
    }

    // 画像変更処理
    if (propsForImageChanging.isImageChanged) {
      const selectedImage = propsForImageChanging.e.target.files[0];
      if (!selectedImage) return;

      const reader = new FileReader();
      reader.onload = (event) => {
        // 画像を表示
        imageElement.src = event.target.result;

        // バーの初期位置を設定
        draggableBarYUpper.style.top = "50%"; // 画像の中央に配置
        draggableBarYUpper.style.transform = "translateY(-50%)";
        draggableBarXUpper.style.left = "0px"; // 画像の左端に配置

        // 座標表示を初期化
        coordinateXUpperDisplay.value = "0";
        coordinateYUpperDisplay.value = `${Math.round(imageElement.offsetHeight / 2)}`;

        // 下部バーも初期化
        draggableBarYLower.style.top = "75%"; // 画像の下部に配置
        draggableBarYLower.style.transform = "translateY(-50%)";
        draggableBarXLower.style.left = "0px"; // 画像の左端に配置

        // 下部座標表示を初期化
        coordinateXLowerDisplay.value = "0";
        coordinateYLowerDisplay.value = `${Math.round(imageElement.offsetHeight * 0.75)}`;
      };

      reader.readAsDataURL(selectedImage);
    }

    console.log("バーの状態更新が完了しました");
  };
};

// 状態更新関数の初期化
const update = updateStateOfBarActivation(0);
