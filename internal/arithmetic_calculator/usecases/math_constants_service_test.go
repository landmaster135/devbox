package usecases

import (
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// #==============================================================#
// ##          MathConstantsService Tests                        ##
// #==============================================================#

// TestNewMathConstantsService は NewMathConstantsService 関数をテストします
func TestNewMathConstantsService(t *testing.T) {
	service := NewMathConstantsService()
	assert.NotNil(t, service, "MathConstantsServiceが正しく作成されませんでした")
}

// TestMathConstantsServiceGetConstants は getConstants メソッドをテストします
func TestMathConstantsServiceGetConstants(t *testing.T) {
	service := NewMathConstantsService()

	constants := service.getConstants()

	// 期待される定数が含まれていることを確認
	expectedConstants := map[string]float64{
		"pi":  math.Pi,
		"e":   math.E,
		"tau": 2 * math.Pi,
	}

	assert.Equal(t, len(expectedConstants), len(constants), "定数の数が期待値と一致しません")

	for name, expectedValue := range expectedConstants {
		actualValue, exists := constants[name]
		assert.True(t, exists, "定数 %s が存在しません", name)
		assert.Equal(t, expectedValue, actualValue, "定数 %s の値が期待値と一致しません", name)
	}
}

// TestMathConstantsServiceGetConstantsValues は各定数の値をテストします
func TestMathConstantsServiceGetConstantsValues(t *testing.T) {
	service := NewMathConstantsService()
	constants := service.getConstants()

	// π (pi) の値をテスト
	pi, exists := constants["pi"]
	assert.True(t, exists, "π定数が存在しません")
	assert.InDelta(t, 3.141592653589793, pi, 1e-15, "π定数の値が正しくありません")

	// e (自然対数の底) の値をテスト
	e, exists := constants["e"]
	assert.True(t, exists, "e定数が存在しません")
	assert.InDelta(t, 2.718281828459045, e, 1e-15, "e定数の値が正しくありません")

	// τ (tau = 2π) の値をテスト
	tau, exists := constants["tau"]
	assert.True(t, exists, "τ定数が存在しません")
	assert.InDelta(t, 6.283185307179586, tau, 1e-15, "τ定数の値が正しくありません")
	assert.InDelta(t, 2*math.Pi, tau, 1e-15, "τ定数は2πと等しくなければなりません")
}

// TestMathConstantsServiceHandleToGetConstants は HandleToGetConstants ハンドラーをテストします
func TestMathConstantsServiceHandleToGetConstants(t *testing.T) {
	service := NewMathConstantsService()

	result, err := service.HandleToGetConstants()

	// エラーが発生しないことを確認
	assert.NoError(t, err, "HandleToGetConstantsでエラーが発生すべきではありません")

	// 結果が空でないことを確認
	assert.NotEmpty(t, result, "結果が空であってはいけません")

	// 結果に期待される文字列が含まれていることを確認
	assert.Contains(t, result, "利用可能な数学定数:", "結果にヘッダーが含まれていません")
	assert.Contains(t, result, "pi =", "結果にπ定数が含まれていません")
	assert.Contains(t, result, "e =", "結果にe定数が含まれていません")
	assert.Contains(t, result, "tau =", "結果にτ定数が含まれていません")

	// 各定数の値が正しく表示されていることを確認
	lines := strings.Split(result, "\n")
	constantsFound := make(map[string]bool)

	for _, line := range lines {
		if strings.Contains(line, "pi =") {
			constantsFound["pi"] = true
			assert.Contains(t, line, "3.141593", "π定数の値が正しく表示されていません")
		}
		if strings.Contains(line, "e =") {
			constantsFound["e"] = true
			assert.Contains(t, line, "2.718282", "e定数の値が正しく表示されていません")
		}
		if strings.Contains(line, "tau =") {
			constantsFound["tau"] = true
			assert.Contains(t, line, "6.283185", "τ定数の値が正しく表示されていません")
		}
	}

	// すべての定数が見つかったことを確認
	assert.True(t, constantsFound["pi"], "π定数が結果に含まれていません")
	assert.True(t, constantsFound["e"], "e定数が結果に含まれていません")
	assert.True(t, constantsFound["tau"], "τ定数が結果に含まれていません")
}

// TestMathConstantsServiceResultFormat は結果のフォーマットをテストします
func TestMathConstantsServiceResultFormat(t *testing.T) {
	service := NewMathConstantsService()

	result, err := service.HandleToGetConstants()
	assert.NoError(t, err)

	// 結果が改行で終わることを確認
	assert.True(t, strings.HasSuffix(result, "\n"), "結果は改行で終わるべきです")

	// 各行が適切なフォーマットであることを確認
	lines := strings.Split(strings.TrimSpace(result), "\n")
	assert.True(t, len(lines) >= 4, "結果は少なくとも4行（ヘッダー + 3つの定数）であるべきです")

	// ヘッダー行の確認
	assert.Equal(t, "利用可能な数学定数:", lines[0], "ヘッダー行が正しくありません")

	// 定数行のフォーマット確認
	constantLines := lines[1:]
	for _, line := range constantLines {
		if strings.TrimSpace(line) != "" {
			assert.Contains(t, line, " = ", "定数行は ' = ' を含むべきです")
			parts := strings.Split(line, " = ")
			assert.Equal(t, 2, len(parts), "定数行は名前と値の2つの部分に分かれるべきです")
		}
	}
}

// TestMathConstantsServiceConstantsImmutability は定数の不変性をテストします
func TestMathConstantsServiceConstantsImmutability(t *testing.T) {
	service := NewMathConstantsService()

	// 複数回呼び出して同じ結果が得られることを確認
	constants1 := service.getConstants()
	constants2 := service.getConstants()

	assert.Equal(t, len(constants1), len(constants2), "複数回の呼び出しで定数の数が変わってはいけません")

	for name, value1 := range constants1 {
		value2, exists := constants2[name]
		assert.True(t, exists, "定数 %s が2回目の呼び出しで存在しません", name)
		assert.Equal(t, value1, value2, "定数 %s の値が複数回の呼び出しで変わってはいけません", name)
	}
}

// TestMathConstantsServiceMathLibraryConsistency は標準ライブラリとの一貫性をテストします
func TestMathConstantsServiceMathLibraryConsistency(t *testing.T) {
	service := NewMathConstantsService()
	constants := service.getConstants()

	// math.Piとの一貫性
	pi, exists := constants["pi"]
	assert.True(t, exists)
	assert.Equal(t, math.Pi, pi, "π定数はmath.Piと一致するべきです")

	// math.Eとの一貫性
	e, exists := constants["e"]
	assert.True(t, exists)
	assert.Equal(t, math.E, e, "e定数はmath.Eと一致するべきです")

	// τ = 2πの関係性
	tau, exists := constants["tau"]
	assert.True(t, exists)
	assert.Equal(t, 2*math.Pi, tau, "τ定数は2πと一致するべきです")
}

// TestMathConstantsServiceEdgeCases は MathConstantsService の境界値ケースをテストします
func TestMathConstantsServiceEdgeCases(t *testing.T) {
	service := NewMathConstantsService()

	// 複数回のHandleToGetConstants呼び出し
	t.Run("複数回のハンドラー呼び出し", func(t *testing.T) {
		result1, err1 := service.HandleToGetConstants()
		result2, err2 := service.HandleToGetConstants()

		assert.NoError(t, err1)
		assert.NoError(t, err2)
		// 結果の長さが同じことを確認（順序は保証されない）
		assert.Equal(t, len(result1), len(result2), "複数回の呼び出しで結果の長さが変わってはいけません")
		assert.Contains(t, result1, "pi =", "結果にπ定数が含まれているべきです")
		assert.Contains(t, result2, "pi =", "結果にπ定数が含まれているべきです")
	})

	// 定数の精度確認
	t.Run("定数の精度確認", func(t *testing.T) {
		constants := service.getConstants()

		// πの精度（少なくとも15桁の精度）
		pi := constants["pi"]
		assert.InDelta(t, math.Pi, pi, 1e-15, "π定数の精度が不十分です")

		// eの精度（少なくとも15桁の精度）
		e := constants["e"]
		assert.InDelta(t, math.E, e, 1e-15, "e定数の精度が不十分です")

		// τの精度（少なくとも15桁の精度）
		tau := constants["tau"]
		assert.InDelta(t, 2*math.Pi, tau, 1e-15, "τ定数の精度が不十分です")
	})
}

// TestMathConstantsServiceStringBuilderUsage は strings.Builder の使用をテストします
func TestMathConstantsServiceStringBuilderUsage(t *testing.T) {
	service := NewMathConstantsService()

	// 結果の構造を確認
	result, err := service.HandleToGetConstants()
	assert.NoError(t, err)

	// 結果が適切に構築されていることを確認
	lines := strings.Split(result, "\n")

	// 空行を除いた有効な行数を確認
	validLines := 0
	for _, line := range lines {
		if strings.TrimSpace(line) != "" {
			validLines++
		}
	}

	// ヘッダー + 3つの定数 = 4行
	assert.Equal(t, 4, validLines, "有効な行数が期待値と一致しません")

	// 各定数行が適切なフォーマットであることを再確認
	constantCount := 0
	for _, line := range lines {
		if strings.Contains(line, " = ") {
			constantCount++
		}
	}
	assert.Equal(t, 3, constantCount, "定数の行数が期待値と一致しません")
}
