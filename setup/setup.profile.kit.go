package setup

import (
	"context"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/go-kratos/kratos/v2/log"
	"go.uber.org/fx"

	"github.com/eden-quan/go-biz-kit/config/def"
)

func NewProfile(lifecycle fx.Lifecycle, conf *def.Configuration, logger log.Logger) {

	var helper = log.NewHelper(logger)
	var fCpu *os.File = nil
	var fMem *os.File = nil
	var err error = nil

	lifecycle.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if !conf.Profile.EnableCpu {
				return nil
			}
			name := conf.Profile.CpuFile
			if name == "" {
				helper.Warn("cpu profile is enable but cpu output file is empty, use cpu.prof")
				name = "cpu.prof"
			}

			fCpu, err = os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
			if err != nil {
				log.Errorf("create cpu profile file %s with error %s", name, err)
				return err
			}

			err = pprof.StartCPUProfile(fCpu)
			if err != nil {
				log.Errorf("start cpu profile with error %s", err)
				return err
			}

			return nil
		},
		OnStop: func(_ context.Context) error {

			if conf.Profile.EnableCpu {
				pprof.StopCPUProfile()
				_ = fCpu.Close()
			}

			if !conf.Profile.EnableMem {
				return nil
			}
			fmt.Println("stop mem")

			name := conf.Profile.MemFile
			if name == "" {
				helper.Warn("mem profile is enable but mem output file is empty, use mem.prof")
				name = "mem.prof"
			}

			fMem, err = os.OpenFile(name, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0666)
			if err != nil {
				helper.Errorf("create memory profile file %s failed with error %s", name, err)
				return err
			}

			err = pprof.WriteHeapProfile(fMem)
			if err != nil {
				helper.Errorf("dumps mem profile to file %s with error %s", name, err)
			}

			_ = fMem.Close()
			return err

		},
	})
}
