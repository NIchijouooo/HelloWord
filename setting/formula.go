package setting

import (
	"fmt"
	"github.com/dengsgo/math-engine/engine"
)

func FormulaTest() {

	s := "1+1"
	r, err := engine.ParseAndExec(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s = %v\r\n", s, r)
}

func FormulaRun(str string) (error, float64) {

	r, err := engine.ParseAndExec(str)
	if err != nil {
		ZAPS.Errorf("公式执行错误 %v", err)
		return err, 0
	}
	//ZAPS.Debugf("公式执行成功，%s = %v", str, r)

	return nil, r
}
