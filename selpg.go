package main
import (
    "flag"
    "fmt"
    "os"
    "bufio"
    "os/exec"
    "log"
    "io"
)
//types相当于结构体
type selpg struct {
    start_page int
    end_page int
    in_filename string
    page_len int
    page_type int/* 'l' for lines-delimited, 'f' for form-feed-delimited */
					/* default is 'l' */
    print_dest string
}
func usage() {
	fmt.Fprintf(os.Stderr,
		`usage: [-s start page(>=1)] [e end page(>=s)] [-l length of page(default 72)] [-f type of file(default 1)] [-d dest] [filename specify input file]
`)
}
var progName string
//main函数
func main() {
    //初始化sa
    sa := selpg {}
    process_args(&sa)
    process_input(sa)
}
// 用flag对参数进行处理
func process_args(sa *selpg) {
	args := os.Args
	progName = args[0]

	if len(args) < 3 {
		log.Printf("%s: not enough arguments\n", progName)
		flag.Usage()
		os.Exit(1)
	}
  var fileType bool

  flag.IntVar(&sa.start_page,"s", -1, "起始页(>1)")
  flag.IntVar(&sa.end_page,"e", -1, "结束页(>=start_page)")
  flag.IntVar(&sa.page_len,"l", 72, "每页的长度")
  flag.BoolVar(&fileType,"f", false, "-f根据换页符确定页数，不能与page_len同时存在（两种情况")
  flag.StringVar(&sa.print_dest,"d", "", "打印的目的地名称")

  flag.Parse()
  // 对命令行参数的处理
  if sa.start_page == -1 || sa.end_page == -1 || sa.start_page > sa.end_page || sa.start_page < 1 || sa.end_page < 1 {
    flag.Usage()
    return
  }
  if fileType {
		sa.page_type = 'f'
	} else {
		sa.page_type = 'l'
	}
  if len(flag.Args()) > 1 {
    flag.Usage()
    return
  }
  if len(flag.Args()) == 1 {
    sa.in_filename = flag.Args()[0]
  }
}

// 对文件进行处理
func process_input(sa selpg) {
	var fin *os.File
	if len(sa.in_filename) == 0 {
		fin = os.Stdin
	} else {
		var err error
		fin, err = os.Open(sa.in_filename)
		if err != nil {
			log.Printf("%s: can't open file\n", progName)
			os.Exit(2)
		}
		defer fin.Close()
	}

	bufFin := bufio.NewReader(fin)

	var fout io.WriteCloser
	cmd := &exec.Cmd{}
//设置输出的目的地
	if len(sa.print_dest) == 0 {
		fout = os.Stdout
	} else {
		cmd = exec.Command("lp", "-d", sa.print_dest)
		cmd.Stdout = os.Stdout
		var err error
		fout, err = cmd.StdinPipe()
		if err != nil {
			log.Printf("%s: can't open pipe\n", progName)
			os.Exit(3)
		}

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
			os.Exit(4)
		}
	}
//按行读取，对应将c语言翻译即可
	if sa.page_type == 'l' {
		line_ctr := 0
		page_ctr := 1
		for {
			line, crc := bufFin.ReadString('\n')
			if crc != nil {
				break
			}
			line_ctr++
			if line_ctr > sa.page_len {
				page_ctr++
				line_ctr = 1
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				_, err := fout.Write([]byte(line))
        if err != nil {
				}
			}
		}
	} else {   //按页读取
		page_ctr := 1
		for {
			line, crc := bufFin.ReadString('\f')
			if crc != nil {
				break
			}
			if (page_ctr >= sa.start_page) && (page_ctr <= sa.end_page) {
				_, err := fout.Write([]byte(line))
        if err != nil {
				}
			}
			page_ctr++
		}
	}
  if err := cmd.Wait(); err != nil {
		//handle err
	}
	fout.Close()
}
