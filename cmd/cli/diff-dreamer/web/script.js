// グローバル変数
let diffResult = null;

// 初期化
document.addEventListener('DOMContentLoaded', function() {
  // 設定が存在する場合は初期値を設定
  if (typeof CONFIG !== 'undefined') {
    if (CONFIG.leftText) {
      document.getElementById('leftText').value = CONFIG.leftText;
    }
    if (CONFIG.rightText) {
      document.getElementById('rightText').value = CONFIG.rightText;
    }
  }
});

// クリア機能
function clearLeft() {
  document.getElementById('leftText').value = '';
  hideResults();
}

function clearRight() {
  document.getElementById('rightText').value = '';
  hideResults();
}

function clearAll() {
  document.getElementById('leftText').value = '';
  document.getElementById('rightText').value = '';
  hideResults();
}

function hideResults() {
  document.getElementById('stats').style.display = 'none';
  document.getElementById('resultSection').style.display = 'none';
}

// 差分計算のメイン関数
function performDiff() {
  const leftText = document.getElementById('leftText').value;
  const rightText = document.getElementById('rightText').value;

  if (!leftText && !rightText) {
    alert('比較するテキストを入力してください。');
    return;
  }

  // 行に分割
  const leftLines = leftText.split('\n');
  const rightLines = rightText.split('\n');

  // 差分計算
  diffResult = calculateDiff(leftLines, rightLines);

  // 結果表示
  displayDiff(diffResult);
  displayStats(diffResult);

  // 結果セクションを表示
  document.getElementById('stats').style.display = 'flex';
  document.getElementById('resultSection').style.display = 'block';
}

// LCS（Longest Common Subsequence）ベースの差分計算
function calculateDiff(leftLines, rightLines) {
  const lcs = computeLCS(leftLines, rightLines);
  const leftResult = [];
  const rightResult = [];

  let leftIndex = 0;
  let rightIndex = 0;
  let lcsIndex = 0;

  while (leftIndex < leftLines.length || rightIndex < rightLines.length) {
    if (
      lcsIndex < lcs.length &&
      leftIndex < leftLines.length &&
      rightIndex < rightLines.length &&
      leftLines[leftIndex] === lcs[lcsIndex] &&
      rightLines[rightIndex] === lcs[lcsIndex]
    ) {
      // 共通行
      leftResult.push({
        content: leftLines[leftIndex],
        type: 'unchanged',
        lineNum: leftIndex + 1
      });
      rightResult.push({
        content: rightLines[rightIndex],
        type: 'unchanged',
        lineNum: rightIndex + 1
      });
      leftIndex++;
      rightIndex++;
      lcsIndex++;
    } else if (
      lcsIndex < lcs.length &&
      leftIndex < leftLines.length &&
      leftLines[leftIndex] === lcs[lcsIndex]
    ) {
      // 右側に追加された行
      rightResult.push({
        content: rightLines[rightIndex],
        type: 'added',
        lineNum: rightIndex + 1
      });
      leftResult.push({
        content: '',
        type: 'empty',
        lineNum: null
      });
      rightIndex++;
    } else if (
      lcsIndex < lcs.length &&
      rightIndex < rightLines.length &&
      rightLines[rightIndex] === lcs[lcsIndex]
    ) {
      // 左側から削除された行
      leftResult.push({
        content: leftLines[leftIndex],
        type: 'deleted',
        lineNum: leftIndex + 1
      });
      rightResult.push({
        content: '',
        type: 'empty',
        lineNum: null
      });
      leftIndex++;
    } else {
      // 変更された行または追加/削除
      if (leftIndex < leftLines.length && rightIndex < rightLines.length) {
        // 変更された行
        const charDiff = calculateCharDiff(leftLines[leftIndex], rightLines[rightIndex]);
        leftResult.push({
          content: leftLines[leftIndex],
          type: 'modified',
          lineNum: leftIndex + 1,
          charDiff: charDiff.left
        });
        rightResult.push({
          content: rightLines[rightIndex],
          type: 'modified',
          lineNum: rightIndex + 1,
          charDiff: charDiff.right
        });
        leftIndex++;
        rightIndex++;
      } else if (leftIndex < leftLines.length) {
        // 削除された行
        leftResult.push({
          content: leftLines[leftIndex],
          type: 'deleted',
          lineNum: leftIndex + 1
        });
        rightResult.push({
          content: '',
          type: 'empty',
          lineNum: null
        });
        leftIndex++;
      } else {
        // 追加された行
        rightResult.push({
          content: rightLines[rightIndex],
          type: 'added',
          lineNum: rightIndex + 1
        });
        leftResult.push({
          content: '',
          type: 'empty',
          lineNum: null
        });
        rightIndex++;
      }
    }
  }

  return {
      left: leftResult,
      right: rightResult
  };
}

// LCS計算
function computeLCS(arr1, arr2) {
  const m = arr1.length;
  const n = arr2.length;
  const dp = Array(m + 1).fill().map(() => Array(n + 1).fill(0));

  // DP テーブル構築
  for (let i = 1; i <= m; i++) {
    for (let j = 1; j <= n; j++) {
      if (arr1[i - 1] === arr2[j - 1]) {
        dp[i][j] = dp[i - 1][j - 1] + 1;
      } else {
        dp[i][j] = Math.max(dp[i - 1][j], dp[i][j - 1]);
      }
    }
  }

  // LCS復元
  const lcs = [];
  let i = m, j = n;
  while (i > 0 && j > 0) {
    if (arr1[i - 1] === arr2[j - 1]) {
      lcs.unshift(arr1[i - 1]);
      i--;
      j--;
    } else if (dp[i - 1][j] > dp[i][j - 1]) {
      i--;
    } else {
      j--;
    }
  }

  return lcs;
}

// 文字レベルの差分計算
function calculateCharDiff(str1, str2) {
  const chars1 = str1.split('');
  const chars2 = str2.split('');
  const lcs = computeLCS(chars1, chars2);

  let result1 = '';
  let result2 = '';
  let i = 0, j = 0, k = 0;

  while (i < chars1.length || j < chars2.length) {
    if (k < lcs.length &&
      i < chars1.length &&
      j < chars2.length &&
      chars1[i] === lcs[k] &&
      chars2[j] === lcs[k]
    ) {
      result1 += chars1[i];
      result2 += chars2[j];
      i++;
      j++;
      k++;
    } else if (
      i < chars1.length &&
      (j >= chars2.length || chars1[i] !== chars2[j])
    ) {
      result1 += `<span class="char-deleted">${escapeHtml(chars1[i])}</span>`;
      i++;
    } else if (j < chars2.length) {
      result2 += `<span class="char-added">${escapeHtml(chars2[j])}</span>`;
      j++;
    }
  }

  return {
      left: result1,
      right: result2
  };
}

// HTML エスケープ
function escapeHtml(text) {
  const div = document.createElement('div');
  div.textContent = text;
  return div.innerHTML;
}

// 差分結果の表示
function displayDiff(result) {
  const leftDiff = document.getElementById('leftDiff');
  const rightDiff = document.getElementById('rightDiff');

  leftDiff.innerHTML = '';
  rightDiff.innerHTML = '';

  for (let i = 0; i < Math.max(result.left.length, result.right.length); i++) {
    const leftLine = result.left[i] || { content: '', type: 'empty', lineNum: null };
    const rightLine = result.right[i] || { content: '', type: 'empty', lineNum: null };

    leftDiff.appendChild(createDiffLine(leftLine));
    rightDiff.appendChild(createDiffLine(rightLine));
  }
}

// 差分行の要素作成
function createDiffLine(line) {
  const div = document.createElement('div');
  div.className = `diff-line ${line.type}`;

  const lineNumber = document.createElement('div');
  lineNumber.className = 'line-number';
  lineNumber.textContent = line.lineNum || '';

  const lineContent = document.createElement('div');
  lineContent.className = 'line-content';

  if (line.type === 'empty') {
    lineContent.innerHTML = '<span class="empty-line">(空行)</span>';
  } else if (line.charDiff) {
    lineContent.innerHTML = line.charDiff;
  } else {
    lineContent.textContent = line.content;
  }

  div.appendChild(lineNumber);
  div.appendChild(lineContent);

  return div;
}

// 統計情報の表示
function displayStats(result) {
  let added = 0, deleted = 0, modified = 0, unchanged = 0;

  result.left.forEach(line => {
    switch (line.type) {
      case 'added': added++; break;
      case 'deleted': deleted++; break;
      case 'modified': modified++; break;
      case 'unchanged': unchanged++; break;
    }
  });

  result.right.forEach(line => {
    switch (line.type) {
      case 'added': added++; break;
      case 'deleted': deleted++; break;
      case 'modified': modified++; break;
      case 'unchanged': unchanged++; break;
    }
  });

  // 重複カウントを修正
  added = result.right.filter(line => line.type === 'added').length;
  deleted = result.left.filter(line => line.type === 'deleted').length;
  modified = result.left.filter(line => line.type === 'modified').length;
  unchanged = result.left.filter(line => line.type === 'unchanged').length;

  document.getElementById('addedLines').textContent = added;
  document.getElementById('deletedLines').textContent = deleted;
  document.getElementById('modifiedLines').textContent = modified;
  document.getElementById('unchangedLines').textContent = unchanged;
}

// キーボードショートカット
document.addEventListener('keydown', function(e) {
  // Ctrl+Enter で比較実行
  if (e.ctrlKey && e.key === 'Enter') {
    e.preventDefault();
    performDiff();
  }

  // Ctrl+Shift+C ですべてクリア
  if (e.ctrlKey && e.shiftKey && e.key === 'C') {
    e.preventDefault();
    clearAll();
  }
});

// テキストエリアの自動リサイズ
function autoResize(textarea) {
  textarea.style.height = 'auto';
  textarea.style.height = textarea.scrollHeight + 'px';
}

// テキストエリアにイベントリスナーを追加
document.addEventListener('DOMContentLoaded', function() {
  const textareas = document.querySelectorAll('textarea');
  textareas.forEach(textarea => {
    textarea.addEventListener('input', function() {
      autoResize(this);
    });
  });
});
