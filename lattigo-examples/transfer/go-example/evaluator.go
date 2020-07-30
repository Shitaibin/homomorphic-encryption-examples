package main

import "github.com/ldsec/lattigo/bfv"

// 本同态加密样例采用的默认参数
var defaultParams = bfv.DefaultParams[bfv.PN13QP218]

// 同态加密执行器
var evaluator bfv.Evaluator

func init() {
	defaultParams.T = 0x3ee0001

	evaluator = bfv.NewEvaluator(defaultParams)
}
