// Copyright (C) 2019-2020 Arnaud Poret
// This work is licensed under the GNU General Public License.
// To view a copy of this license, visit https://www.gnu.org/licenses/gpl.html.
package main
import (
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)
func Stream() {
    var (
        err error
        help,usage,getTerminal bool
        outFile,outFilePath,outFileBase,blackFile string
        args,nodes,blackNodes,seeds,termNodes []string
        edges,ward [][]string
        nodeSP map[string][]string
        edgeNames map[string]map[string][]string
        edgeSP map[string]map[string][][]string
        flagSet *flag.FlagSet
    )
    flagSet=flag.NewFlagSet("",flag.ContinueOnError)
    flagSet.Usage=func() {}
    flagSet.BoolVar(&help,"help",false,"")
    flagSet.BoolVar(&help,"h",false,"")
    flagSet.BoolVar(&usage,"usage",false,"")
    flagSet.BoolVar(&usage,"u",false,"")
    flagSet.StringVar(&outFile,"out","out.sif","")
    flagSet.StringVar(&outFile,"o","out.sif","")
    flagSet.BoolVar(&getTerminal,"terminal",false,"")
    flagSet.BoolVar(&getTerminal,"t",false,"")
    flagSet.StringVar(&blackFile,"blacklist","","")
    flagSet.StringVar(&blackFile,"b","","")
    err=flagSet.Parse(os.Args[2:])
    if err!=nil {
        fmt.Println("Error: pathrider stream: "+err.Error())
    } else if help {
        fmt.Println(strings.Join([]string{
            "",
            "Find the upstream/downstream paths starting from some nodes of interest (the",
            "seed nodes) in a network.",
            "",
            "Typical use is to find in a network the paths regulating some nodes (the",
            "upstream paths) or the paths regulated by some nodes (the downstream paths).",
            "",
            "Usage: pathrider stream [options] <networkFile> <seedFile> <direction>",
            "",
            "Positional arguments:",
            "    * <networkFile>: the network encoded in a SIF file",
            "    * <seedFile>: the seed nodes listed in a file (one node per line)",
            "    * <direction>: follow the up stream (up) or the down stream (down)",
            "",
            "Options:",
            "    * -t/-terminal: also find the terminal nodes reachable from the seed nodes,",
            "                    namely the nodes having no predecessors in case of",
            "                    upstreaming, or the nodes having no successors in case of",
            "                    downstreaming (default: not used by default)",
            "    * -b/-blacklist <file>: a file containing a list of nodes to be blacklisted",
            "                            (one node per line), the paths containing such nodes",
            "                            will not be considered (default: not used by",
            "                            default)",
            "    * -o/-out <file>: the output SIF file (default: out.sif)",
            "    * -u/-usage: print usage only",
            "    * -h/-help: print help",
            "",
            "Output file(s) (unless changed with -o/-out):",
            "    * out.sif: a SIF file encoding the upstream/downstream paths starting from",
            "               the seed nodes in the network",
            "    * out-terminal.txt: a file listing the upstream/downstream terminal nodes",
            "                        reachable from the seed nodes in the network",
            "                        (requires -t/-terminal)",
            "",
            "Cautions:",
            "    * the network must be in the SIF file format (see the readme file of",
            "      pathrider)",
            "    * edge duplicates are automatically removed",
            "    * edges are assumed to be directed",
            "",
            "For more information, see https://github.com/arnaudporet/pathrider.",
            "",
        },"\n"))
    } else if usage {
        fmt.Println(strings.Join([]string{
            "",
            "Usage: pathrider stream [options] <networkFile> <seedFile> <direction>",
            "",
            "Positional arguments:",
            "    * <networkFile>: the network encoded in a SIF file",
            "    * <seedFile>: the seed nodes listed in a file (one node per line)",
            "    * <direction>: follow the up stream (up) or the down stream (down)",
            "",
            "Options:",
            "    * -t/-terminal: also find the terminal nodes reachable from the seed nodes,",
            "                    namely the nodes having no predecessors in case of",
            "                    upstreaming, or the nodes having no successors in case of",
            "                    downstreaming (default: not used by default)",
            "    * -b/-blacklist <file>: a file containing a list of nodes to be blacklisted",
            "                            (one node per line), the paths containing such nodes",
            "                            will not be considered (default: not used by",
            "                            default)",
            "    * -o/-out <file>: the output SIF file (default: out.sif)",
            "    * -u/-usage: print usage only",
            "    * -h/-help: print help",
            "",
            "Output file(s) (unless changed with -o/-out):",
            "    * out.sif: a SIF file encoding the upstream/downstream paths starting from",
            "               the seed nodes in the network",
            "    * out-terminal.txt: a file listing the upstream/downstream terminal nodes",
            "                        reachable from the seed nodes in the network",
            "                        (requires -t/-terminal)",
            "",
        },"\n"))
    } else if filepath.Ext(outFile)!=".sif" {
        fmt.Println("Error: pathrider stream: "+outFile+": must have the \".sif\" file extension")
    } else if len(flagSet.Args())!=3 {
        fmt.Println("Error: pathrider stream: wrong number of positional arguments, expecting: <networkFile> <seedFile> <direction>")
    } else if (flagSet.Arg(2)!="up") && (flagSet.Arg(2)!="down") {
        fmt.Println("Error: pathrider stream: "+flagSet.Arg(2)+": unknown direction, expecting one of: up, down")
    } else {
        args=flagSet.Args()
        fmt.Println("reading network: "+args[0])
        nodes,edges,edgeNames,err=ReadNetwork(args[0])
        if err!=nil {
            fmt.Println("Error: pathrider stream: "+args[0]+": "+err.Error())
        } else {
            if blackFile!="" {
                fmt.Println("reading blacklist: "+blackFile)
                blackNodes,err=ReadNodes(blackFile,nodes)
                if err!=nil {
                    fmt.Println("Error: pathrider stream: "+blackFile+": "+err.Error())
                } else {
                    fmt.Println("blacklisting nodes")
                    nodes,edges,edgeNames,err=RmNodes(edges,edgeNames,blackNodes)
                    if err!=nil {
                        fmt.Println("Error: pathrider stream: "+blackFile+": "+err.Error())
                    }
                }
            }
            if err==nil {
                fmt.Println("reading seed nodes: "+args[1])
                seeds,err=ReadNodes(args[1],nodes)
                if err!=nil {
                    fmt.Println("Error: pathrider stream: "+args[1]+": "+err.Error())
                } else {
                    fmt.Println(args[2]+"streaming seed nodes")
                    if args[2]=="up" {
                        nodeSP,edgeSP=GetPredecessors(edges)
                        ward=BackwardEdges(seeds,nodeSP,edgeSP)
                    } else if args[2]=="down" {
                        nodeSP,edgeSP=GetSuccessors(edges)
                        ward=ForwardEdges(seeds,nodeSP,edgeSP)
                    }
                    if len(ward)==0 {
                        fmt.Println("Warning: pathrider stream: "+args[1]+": no "+args[2]+"stream paths found")
                    } else {
                        fmt.Println("writing "+args[2]+"stream paths: "+outFile)
                        err=WriteNetwork(outFile,ward,edgeNames)
                        if err!=nil {
                            fmt.Println("Error: pathrider stream: "+outFile+": "+err.Error())
                        } else if getTerminal {
                            fmt.Println("computing "+args[2]+"stream terminal nodes")
                            if args[2]=="up" {
                                nodeSP,edgeSP=GetPredecessors(ward)
                            } else if args[2]=="down" {
                                nodeSP,edgeSP=GetSuccessors(ward)
                            }
                            termNodes=TerminalNodes(nodeSP)
                            if len(termNodes)==0 {
                                fmt.Println("Warning: pathrider stream: "+args[1]+": no "+args[2]+"stream terminal nodes found")
                            } else {
                                outFilePath,outFileBase=filepath.Split(outFile)
                                outFileBase=strings.TrimSuffix(outFileBase,".sif")
                                outFileBase+="-terminal.txt"
                                outFile=filepath.Join(outFilePath,outFileBase)
                                fmt.Println("writing "+args[2]+"stream terminal nodes: "+outFile)
                                err=WriteText(outFile,termNodes)
                                if err!=nil {
                                    fmt.Println("Error: pathrider stream: "+outFile+": "+err.Error())
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
