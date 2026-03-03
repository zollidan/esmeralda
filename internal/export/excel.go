package export

import (
	"github.com/xuri/excelize/v2"
	"github.com/zollidan/esmeralda/internal/stats"
)

type Writer struct {
    file   *excelize.File
    sw     *excelize.StreamWriter
    rowIdx int
}

func NewWriter(filename string) (*Writer, error) {
    f := excelize.NewFile()
    sw, err := f.NewStreamWriter("Sheet1")
    if err != nil {
        return nil, err
    }

    w := &Writer{file: f, sw: sw, rowIdx: 1}

    if err := w.writeHeaders(); err != nil {
        return nil, err
    }

    return w, nil
}

func (w *Writer) writeHeaders() error {
    cell, _ := excelize.CoordinatesToCellName(1, w.rowIdx)
    w.rowIdx++
    return w.sw.SetRow(cell, stats.Headers())
}

func (w *Writer) WriteRow(row stats.Row) error {
	cell, _ := excelize.CoordinatesToCellName(1, w.rowIdx)
	w.rowIdx++
	return w.sw.SetRow(cell, row.ToSlice())
}

func (w *Writer) Save(filename string) error {
    if err := w.sw.Flush(); err != nil {
        return err
    }
    return w.file.SaveAs(filename)
}