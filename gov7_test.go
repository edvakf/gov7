package gov7

import "testing"

func TestExec(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	result, err := v7.Exec("1")
	if err != nil {
		t.Fatal(err.Error())
	}

	n, err := v7.ToNumber(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	if n != 1 {
		t.Fatalf("wrong return number: expected 1, but got %f", n)
	}
}

func TestCreateString(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	str := "abcあいうえお123"

	v7Str := v7.CreateString(str)

	s, err := v7.ToString(v7Str)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != str {
		t.Fatalf("wrong return string: want `%s`, got `%s`", str, s)
	}
}

func TestToJSON(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	str := "abcあいうえお123"
	jstr := "\"" + str + "\""

	j := v7.ToJSON(v7.CreateString(str), 256)

	if j != jstr {
		t.Fatalf("error making JSON: want `%s`, got `%s`", jstr, j)
	}
}

func TestExecFunction(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	result, err := v7.Exec("function myfunc(){return 'str'}")
	if err != nil {
		t.Fatal(err.Error())
	}

	result, err = v7.Exec("myfunc()")
	if err != nil {
		t.Fatal(err.Error())
	}

	s, err := v7.ToString(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != "str" {
		t.Fatalf("wrong return string: want \"str\", got %s", s)
	}
}

func TestApply(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	str := "abcあいうえお123"

	result, err := v7.Exec("function myfunc(x){return x}")
	if err != nil {
		t.Fatal(err.Error())
	}

	global := v7.GetGlobalObject()
	if v7.IsUndefined(global) {
		t.Fatal("could not get global object")
	}

	ary := v7.CreateArray()
	v7.ArrayPush(ary, v7.CreateString(str))

	fn := v7.Get(global, "myfunc")
	if v7.IsUndefined(global) {
		t.Fatal("could not get target function")
	}

	result, err = v7.Apply(fn, global, ary)
	if err != nil {
		t.Fatalf("error calling function: %s", err.Error())
	}

	s, err := v7.ToString(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != str {
		t.Fatalf("wrong return string: want `%s`, got `%s`", str, s)
	}
}

func TestSetGet(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	global := v7.GetGlobalObject()
	if v7.IsUndefined(global) {
		t.Fatal("could not get global object")
	}

	str := "abcあいうえお123"

	v7Str := v7.CreateString(str)

	err := v7.Set(global, "myproperty", PROPERTY_READ_ONLY, v7Str)
	if err != nil {
		t.Fatal("could not set property")
	}

	result := v7.Get(global, "myproperty")
	if v7.IsUndefined(global) {
		t.Fatal("could not get target property")
	}

	s, err := v7.ToString(result)
	if err != nil {
		t.Fatal(err.Error())
	}
	if s != str {
		t.Fatalf("wrong return string: want \"str\", got %s", s)
	}

}

func TestSyntaxError(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	_, err := v7.Exec("var;")
	if err != nil {
		t.Logf("expected syntax error: %s", err.Error())
	} else {
		t.Fatal("syntax error did not occur")
	}
}

func TestJSError(t *testing.T) {
	v7 := New()
	defer v7.Destroy()

	_, err := v7.Exec("throw new Error('yey')")
	if err != nil {
		t.Logf("expected syntax error: %s", err.Error())
	} else {
		t.Fatal("syntax error did not occur")
	}
}
