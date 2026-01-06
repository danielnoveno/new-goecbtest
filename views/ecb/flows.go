/*
   file:           views/ecb/flows.go
   description:    Shared helpers for line flows
*/

package ecb

import (
	"fmt"
	"strings"

	"github.com/go-gorp/gorp"
	"go-ecb/app/types"
	"go-ecb/configs"
	"go-ecb/services/ecbcore"
)

type lineConfig struct {
	lineID     string
	workcenter string
	category   string
}

func parseList(value string) []string {
	if strings.TrimSpace(value) == "" {
		return nil
	}
	parts := strings.Split(value, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		if trimmed := strings.TrimSpace(part); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func valueAt(list []string, idx int) string {
	if len(list) == 0 {
		return ""
	}
	if idx < 0 {
		idx = 0
	}
	if idx >= len(list) {
		idx = 0
	}
	return list[idx]
}

func deriveCategory(lineID string) string {
	if lineID == "" {
		return "REF"
	}
	parts := strings.Fields(lineID)
	if len(parts) == 0 || parts[0] == "" {
		return "REF"
	}
	return parts[0]
}

func buildLineConfig(cfg configs.SimoConfig, lineIndex int) lineConfig {
	listIDs := parseList(cfg.EcbLineIds)
	workcenters := parseList(cfg.EcbWorkcenters)
	lineID := valueAt(listIDs, lineIndex)
	if lineID == "" {
		lineID = valueAt(listIDs, 0)
	}
	workcenter := valueAt(workcenters, lineIndex)
	if workcenter == "" {
		workcenter = valueAt(workcenters, 0)
	}
	return lineConfig{
		lineID:     lineID,
		workcenter: workcenter,
		category:   deriveCategory(lineID),
	}
}

func ConfigureSnOnlyLine(line *LineState, dbmap *gorp.DbMap, cfg configs.SimoConfig, lineIndex int) {
	if line == nil {
		return
	}
	line.SetSuccessMessage("data berhasil disimpan")
	config := buildLineConfig(cfg, lineIndex)
	remote := ecbcore.NewRemoteChecker(configs.LoadConfig())
	coreCfg := ecbcore.LineConfig{LineID: config.lineID, Workcenter: config.workcenter, Category: config.category}
	line.SetSerialValidator(func(input string) (string, error) {
		return ecbcore.ValidateSnOnlySerial(dbmap, remote, input)
	})
	line.SetSaveHandler(func() error {
		sn, fg, _, _, _ := line.Values()
		return ecbcore.SaveSnOnly(dbmap, coreCfg, sn, fg)
	})
}

type refrigFlowContext struct {
	line           *LineState
	dbmap          *gorp.DbMap
	cfg            configs.SimoConfig
	lineIndex      int
	lineCfg        lineConfig
	isPO           bool
	remote         ecbcore.RemoteChecker
	lastCompressor *types.Compressor
	lastEcbPo      *types.EcbPo
	poValue        string
}

func newRefrigFlowContext(line *LineState, dbmap *gorp.DbMap, cfg configs.SimoConfig, lineIndex int, isPO bool) *refrigFlowContext {
	return &refrigFlowContext{
		line:      line,
		dbmap:     dbmap,
		cfg:       cfg,
		lineIndex: lineIndex,
		lineCfg:   buildLineConfig(cfg, lineIndex),
		isPO:      isPO,
		remote:    ecbcore.NewRemoteChecker(configs.LoadConfig()),
	}
}

func ConfigureRefrigFlow(line *LineState, dbmap *gorp.DbMap, cfg configs.SimoConfig, lineIndex int, isPO bool) {
	if line == nil {
		return
	}
	line.SetSuccessMessage("data berhasil disimpan")
	ctx := newRefrigFlowContext(line, dbmap, cfg, lineIndex, isPO)
	line.SetStepValidator(StepSPC, ctx.validateSpc)
	line.SetSerialValidator(ctx.validateSerial)
	line.SetStepValidator(StepCompressorType, ctx.validateCompressorType)
	line.SetStepValidator(StepCompressorCode, ctx.validateCompressorCode)
	line.SetSaveHandler(ctx.saveData)
}

func (ctx *refrigFlowContext) validateSpc(raw string) error {
	return ecbcore.ValidateSpc(ctx.dbmap, ctx.remote, raw)
}

func (ctx *refrigFlowContext) validateSerial(raw string) (string, error) {
	return ecbcore.ValidateRefrigSerial(ctx.dbmap, ctx.remote, raw)
}

func (ctx *refrigFlowContext) validateCompressorType(raw string) error {
	sn, _, _, _, _ := ctx.line.Values()
	comp, po, err := ecbcore.ValidateCompressorType(ctx.dbmap, sn, raw, ctx.isPO)
	if err != nil {
		return err
	}
	ctx.lastCompressor = comp
	ctx.lastEcbPo = po
	if po != nil {
		ctx.poValue = po.Po
	}
	return nil
}

func (ctx *refrigFlowContext) validateCompressorCode(raw string) error {
	return ecbcore.ValidateCompressorCode(ctx.lastCompressor, raw)
}

func (ctx *refrigFlowContext) saveData() error {
	if ctx.line == nil {
		return fmt.Errorf("line tidak tersedia")
	}
	sn, fg, spc, compType, compCode := ctx.line.Values()
	coreCfg := ecbcore.LineConfig{LineID: ctx.lineCfg.lineID, Workcenter: ctx.lineCfg.workcenter, Category: ctx.lineCfg.category}
	var poID *int
	if ctx.lastEcbPo != nil {
		poID = &ctx.lastEcbPo.ID
	}
	if err := ecbcore.SaveRefrig(ctx.dbmap, coreCfg, sn, fg, spc, compType, compCode, ctx.poValue, ctx.isPO, poID); err != nil {
		return err
	}
	ctx.poValue = ""
	ctx.lastCompressor = nil
	ctx.lastEcbPo = nil
	return nil
}
