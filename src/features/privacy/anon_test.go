package privacy

import "testing"

func TestGenerateAnonymousNumber(t *testing.T) {
	result := generateAnonymousNumber(12345)
	if len(result) != 4 {
		t.Errorf("generateAnonymousNumber(12345) = %q, want 4-digit string", result)
	}

	result2 := generateAnonymousNumber(12345)
	if result != result2 {
		t.Errorf("generateAnonymousNumber should be deterministic: %q != %q", result, result2)
	}

	result3 := generateAnonymousNumber(99999)
	if len(result3) != 4 {
		t.Errorf("generateAnonymousNumber(99999) = %q, want 4-digit string", result3)
	}

	resultZero := generateAnonymousNumber(0)
	if len(resultZero) != 4 {
		t.Errorf("generateAnonymousNumber(0) = %q, want 4-digit string", resultZero)
	}

	resultNeg := generateAnonymousNumber(-100)
	if len(resultNeg) != 4 {
		t.Errorf("generateAnonymousNumber(-100) = %q, want 4-digit string", resultNeg)
	}
}
