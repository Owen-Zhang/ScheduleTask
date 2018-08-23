package etc

import 	"flag"

func init() {
	var etcfile string
	flag.StringVar(&etcfile, "f", "etc/worker.yml", "worker etc file.")
	flag.Parse()

	if err := newconfig(etcfile); err != nil {
		panic(err)
	}
}
