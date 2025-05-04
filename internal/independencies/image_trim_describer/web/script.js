// 画像トリミングツール用JavaScript

// DOM要素の取得
const imageSelector = document.getElementById('imageSelector');
const imageLoader = document.getElementById('imageLoader');
const imageElement = document.getElementById('imageElement');
const imageView = document.getElementById('imageView');
const trimButton = document.getElementById('trimButton');
const resetButton = document.getElementById('resetButton');
const resultContainer = document.getElementById('result');
const clipboardCoordinates = document.getElementById('clipboardCoordinates');

// 座標表示要素
const coordinateXUpperDisplay = document.getElementById('coordinateXUpperDisplay');
const coordinateYUpperDisplay = document.getElementById('coordinateYUpperDisplay');
const coordinateXLowerDisplay = document.getElementById('coordinateXLowerDisplay');
const coordinateYLowerDisplay = document.getElementById('coordinateYLowerDisplay');

// ドラッグ可能なバー要素
const draggableBarXUpper = document.getElementById('draggableBarXUpper');
const draggableBarYUpper = document.getElementById('draggableBarYUpper');
const draggableBarXLower = document.getElementById('draggableBarXLower');
const draggableBarYLower = document.getElementById('draggableBarYLower');

// ステータス表示要素
const lockStatus = document.getElementById('lockStatus');
const copyStatus = document.getElementById('copyStatus');

// 状態管理変数
let activeBarIndex = 0; // アクティブなバーのインデックス（0: 上部, 1: 下部）
let isLocked = false; // バーがロックされているかどうか
let selectedFilename = ''; // 選択されたファイル名

// 初期化処理
window.addEventListener('DOMContentLoaded', () => {
  // 設定情報が存在する場合は画像一覧を設定
  if (typeof CONFIG !== 'undefined' && CONFIG.imageFiles) {
    loadImageList(CONFIG.imageFiles);
  }

  // イベントリスナーの設定
  imageSelector.addEventListener('change', handleImageSelect);
  imageLoader.addEventListener('change', handleImageUpload);
  document.addEventListener('mousemove', handleMouseMove);
  document.addEventListener('mousedown', handleMouseDown);
  document.addEventListener('mouseup', handleMouseUp);
  trimButton.addEventListener('click', handleTrimRequest);
  resetButton.addEventListener('click', resetCoordinates);
});

// 画像一覧の読み込み
function loadImageList(images) {
  // セレクトボックスに画像一覧を追加
  images.forEach(image => {
    const option = document.createElement('option');
    option.value = image;
    option.textContent = image;
    imageSelector.appendChild(option);
  });
}

// 画像選択時の処理
function handleImageSelect(event) {
  const filename = event.target.value;
  if (!filename) return;

  selectedFilename = filename;

  // 設定情報が存在する場合はファイルパスを構築
  if (typeof CONFIG !== 'undefined' && CONFIG.srcDir) {
    // ファイルパスを直接指定できないため、Data URLを使用
    loadImageFromFile(filename);
  } else {
    // 従来の方法（サーバーから取得）
    imageElement.src = `/images/${filename}`;
    imageElement.onload = initializeCoordinates;
  }

  enableTrimButton();
}

// ファイルからの画像読み込み
function loadImageFromFile(filename) {
  // imagesディレクトリから画像を読み込む
  imageElement.src = `images/${filename}`;
  imageElement.alt = `選択された画像: ${filename}`;

  // 画像が読み込まれたら座標を初期化
  imageElement.onload = initializeCoordinates;
}

// ファイルアップロード時の処理
function handleImageUpload(event) {
  const file = event.target.files[0];
  if (!file) return;

  const reader = new FileReader();
  reader.onload = function (e) {
    imageElement.src = e.target.result;
    imageElement.onload = initializeCoordinates;
  };
  reader.readAsDataURL(file);

  // アップロードされたファイルはローカルファイルとして扱う
  selectedFilename = file.name;
  enableTrimButton();
}

// 座標の初期化
function initializeCoordinates() {
  // 上部バーの初期位置（左上）
  draggableBarXUpper.style.left = '0px';
  draggableBarYUpper.style.top = '0px';
  coordinateXUpperDisplay.value = '0';
  coordinateYUpperDisplay.value = '0';

  // 下部バーの初期位置（右下）
  const imgWidth = imageElement.offsetWidth || 400;
  const imgHeight = imageElement.offsetHeight || 300;

  draggableBarXLower.style.left = `${imgWidth}px`;
  draggableBarYLower.style.top = `${imgHeight}px`;
  coordinateXLowerDisplay.value = imgWidth.toString();
  coordinateYLowerDisplay.value = imgHeight.toString();

  // バーの表示
  draggableBarXUpper.style.display = 'block';
  draggableBarYUpper.style.display = 'block';
  draggableBarXLower.style.display = 'block';
  draggableBarYLower.style.display = 'block';

  // 結果表示をクリア
  resultContainer.innerHTML = '';
}

// マウス移動時の処理
function handleMouseMove(event) {
  if (isLocked) return;

  const imageRect = imageElement.getBoundingClientRect();

  // 画像の範囲内かどうかを確認
  if (event.clientX >= imageRect.left && event.clientX <= imageRect.right &&
    event.clientY >= imageRect.top && event.clientY <= imageRect.bottom) {

    // 画像内の相対座標を計算
    const x = event.clientX - imageRect.left;
    const y = event.clientY - imageRect.top;

    // アクティブなバーに応じて座標を更新
    if (activeBarIndex === 0) {
      // 上部バー
      draggableBarXUpper.style.left = `${x}px`;
      draggableBarYUpper.style.top = `${y}px`;
      coordinateXUpperDisplay.value = Math.round(x).toString();
      coordinateYUpperDisplay.value = Math.round(y).toString();
    } else {
      // 下部バー
      draggableBarXLower.style.left = `${x}px`;
      draggableBarYLower.style.top = `${y}px`;
      coordinateXLowerDisplay.value = Math.round(x).toString();
      coordinateYLowerDisplay.value = Math.round(y).toString();
    }

    // トリミングボタンの有効化チェック
    checkCoordinatesValidity();
  }
}

// マウスダウン時の処理
function handleMouseDown(event) {
  // ロック状態を切り替え
  isLocked = true;
  lockStatus.textContent = 'ロック中';
  lockStatus.style.color = 'red';

  // 座標をクリップボードにコピー
  copyCoordinatesToClipboard();
}

// マウスアップ時の処理
function handleMouseUp(event) {
  // ロック状態を解除
  isLocked = false;
  lockStatus.textContent = '移動可能';
  lockStatus.style.color = 'green';

  // アクティブなバーを切り替え
  activeBarIndex = (activeBarIndex + 1) % 2;
}

// 座標をクリップボードにコピー
function copyCoordinatesToClipboard() {
  const x1 = coordinateXUpperDisplay.value;
  const y1 = coordinateYUpperDisplay.value;
  const x2 = coordinateXLowerDisplay.value;
  const y2 = coordinateYLowerDisplay.value;

  const copyText = `${x1} ${y1} ${x2} ${y2}`;

  // クリップボード表示用テキストボックスに座標を表示
  clipboardCoordinates.value = copyText;

  navigator.clipboard.writeText(copyText).then(() => {
    copyStatus.textContent = 'コピーしました: ' + copyText;
    setTimeout(() => {
      copyStatus.textContent = '';
    }, 2000);
  }, (err) => {
    console.error('クリップボードへのコピーに失敗しました:', err);
  });
}

// 座標の有効性チェック
function checkCoordinatesValidity() {
  const x1 = parseInt(coordinateXUpperDisplay.value, 10) || 0;
  const y1 = parseInt(coordinateYUpperDisplay.value, 10) || 0;
  const x2 = parseInt(coordinateXLowerDisplay.value, 10) || 0;
  const y2 = parseInt(coordinateYLowerDisplay.value, 10) || 0;

  // 座標が有効かどうかを確認（x2 > x1, y2 > y1）
  const isValid = x2 > x1 && y2 > y1;

  // トリミングボタンの有効/無効を切り替え
  if (isValid && selectedFilename) {
    enableTrimButton();
  } else {
    disableTrimButton();
  }

  return isValid;
}

// トリミングボタンを有効化
function enableTrimButton() {
  trimButton.disabled = false;
}

// トリミングボタンを無効化
function disableTrimButton() {
  trimButton.disabled = true;
}

// 座標のリセット
function resetCoordinates() {
  if (imageElement.src || imageElement.alt) {
    initializeCoordinates();
    // クリップボード表示用テキストボックスをクリア
    clipboardCoordinates.value = '';
  }
}

// トリミングリクエストの処理
function handleTrimRequest() {
  if (!checkCoordinatesValidity() || !selectedFilename) {
    showResult('無効な座標またはファイルが選択されていません', 'error');
    return;
  }

  const x1 = parseInt(coordinateXUpperDisplay.value, 10);
  const y1 = parseInt(coordinateYUpperDisplay.value, 10);
  const x2 = parseInt(coordinateXLowerDisplay.value, 10);
  const y2 = parseInt(coordinateYLowerDisplay.value, 10);

  // 座標情報をクリップボードにコピー
  const coordText = `${x1} ${y1} ${x2} ${y2}`;

  // クリップボード表示用テキストボックスに座標を表示
  clipboardCoordinates.value = coordText;

  navigator.clipboard.writeText(coordText).then(() => {
    showResult(`座標情報をクリップボードにコピーしました: ${coordText}`, 'success');
    showResult(`この座標情報を使用して、コマンドラインから画像トリミングを実行できます:`, 'info');

    // 設定情報が存在する場合はコマンド例を表示
    if (typeof CONFIG !== 'undefined') {
      const commandExample = `image-trimmer -src "${CONFIG.srcDir}" -out "${CONFIG.outDir}" -x1 ${x1} -y1 ${y1} -x2 ${x2} -y2 ${y2} ${CONFIG.move ? '-move' : ''}`;
      showResult(`<code>${commandExample}</code>`, 'command');
    }
  }, (err) => {
    console.error('クリップボードへのコピーに失敗しました:', err);
    showResult('クリップボードへのコピーに失敗しました', 'error');
  });
}

// 結果の表示
function showResult(message, type) {
  const resultElement = document.createElement('p');
  resultElement.className = type;
  resultElement.innerHTML = message;

  // 結果を追加（古い結果は残す）
  resultContainer.appendChild(resultElement);

  // 結果表示エリアまでスクロール
  resultContainer.scrollIntoView({ behavior: 'smooth', block: 'end' });
}
