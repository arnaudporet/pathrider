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
func Connect() {
    var (
        err1,err2 error
        help,usage,getShortest bool
        outFile,outFilePath,outFileBase,blackFile string
        args,nodes,sources,targets,blackNodes,selfLooped []string
        edges,forward,backward,intersect,noSelfLoop,allShortest [][]string
        nodeSucc,nodePred map[string][]string
        edgeNames map[string]map[string][]string
        edgeSucc,edgePred map[string]map[string][][]string
        flagSet *flag.FlagSet
    )
    flagSet=flag.NewFlagSet("",flag.ContinueOnError)
    flagSet.Usage=func() {}
    flagSet.BoolVar(&help,"help",false,"")
    flagSet.BoolVar(&help,"h",false,"")
    flagSet.BoolVar(&usage,"usage",false,"")
    flagSet.BoolVar(&usage,"u",false,"")
    flagSet.BoolVar(&getShortest,"shortest",false,"")
    flagSet.BoolVar(&getShortest,"s",false,"")
    flagSet.StringVar(&outFile,"out","out.sif","")
    flagSet.StringVar(&outFile,"o","out.sif","")
    flagSet.StringVar(&blackFile,"blacklist","","")
    flagSet.StringVar(&blackFile,"b","","")
    err1=flagSet.Parse(os.Args[2:])
    if err1!=nil {
        fmt.Println("Error: pathrider connect: "+err1.Error())
    } else if help {
        fmt.Println(strings.Join([]string{
            "",
            "Find the paths connecting some nodes of interest in a network.",
            "",
            "Typical use is to find in a network the paths connecting some source nodes to",
            "some target nodes.",
            "",
            "Usage: pathrider connect [options] <networkFile> <sourceFile> <targetFile>",
            "",
            "Positional arguments:",
            "    * <networkFile>: the network encoded in a SIF file",
            "    * <sourceFile>: the source nodes listed in a file (one node per line)",
            "    * <targetFile>: the target nodes listed in a file (one node per line)",
            "    * if sources = targets then provide the same node list twice",
            "",
            "Options:",
            "    * -s/-shortest: also find the shortest connecting paths (default: not used",
            "                    by default)",
            "    * -b/-blacklist <file>: a file containing a list of nodes to be blacklisted",
            "                            (one node per line), the paths containing such nodes",
            "                            will not be considered (default: not used by",
            "                            default)",
            "    * -o/-out <file>: the output SIF file (default: out.sif)",
            "    * -u/-usage: print usage only",
            "    * -h/-help: print help",
            "",
            "Output file(s) (unless changed with -o/-out):",
            "    * out.sif: a SIF file encoding all the paths connecting the source nodes to",
            "               the target nodes in the network",
            "    * out-shortest.sif: a SIF file encoding only the shortest connecting paths",
            "                        (requires -s/-shortest)",
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
            "Usage: pathrider connect [options] <networkFile> <sourceFile> <targetFile>",
            "",
            "Positional arguments:",
            "    * <networkFile>: the network encoded in a SIF file",
            "    * <sourceFile>: the source nodes listed in a file (one node per line)",
            "    * <targetFile>: the target nodes listed in a file (one node per line)",
            "    * if sources = targets then provide the same node list twice",
            "",
            "Options:",
            "    * -s/-shortest: also find the shortest connecting paths (default: not used",
            "                    by default)",
            "    * -b/-blacklist <file>: a file containing a list of nodes to be blacklisted",
            "                            (one node per line), the paths containing such nodes",
            "                            will not be considered (default: not used by",
            "                            default)",
            "    * -o/-out <file>: the output SIF file (default: out.sif)",
            "    * -u/-usage: print usage only",
            "    * -h/-help: print help",
            "",
            "Output file(s) (unless changed with -o/-out):",
            "    * out.sif: a SIF file encoding all the paths connecting the source nodes to",
            "               the target nodes in the network",
            "    * out-shortest.sif: a SIF file encoding only the shortest connecting paths",
            "                        (requires -s/-shortest)",
            "",
        },"\n"))
    } else if filepath.Ext(outFile)!=".sif" {
        fmt.Println("Error: pathrider connect: "+outFile+": must have the \".sif\" file extension")
    } else if len(flagSet.Args())!=3 {
        fmt.Println("Error: pathrider connect: wrong number of positional arguments, expecting: <networkFile> <sourceFile> <targetFile>")
    } else {
        args=flagSet.Args()
        fmt.Println("reading network: "+args[0])
        nodes,edges,edgeNames,err1=ReadNetwork(args[0])
        if err1!=nil {
            fmt.Println("Error: pathrider connect: "+args[0]+": "+err1.Error())
        } else {
            if blackFile!="" {
                fmt.Println("reading blacklist: "+blackFile)
                blackNodes,err1=ReadNodes(blackFile,nodes)
                if err1!=nil {
                    fmt.Println("Error: pathrider connect: "+blackFile+": "+err1.Error())
                } else {
                    fmt.Println("blacklisting nodes")
                    nodes,edges,edgeNames,err1=RmNodes(edges,edgeNames,blackNodes)
                    if err1!=nil {
                        fmt.Println("Error: pathrider connect: "+blackFile+": "+err1.Error())
                    }
                }
            }
            if err1==nil {
                fmt.Println("reading source nodes: "+args[1])
                sources,err1=ReadNodes(args[1],nodes)
                fmt.Println("reading target nodes: "+args[2])
                targets,err2=ReadNodes(args[2],nodes)
                if err1!=nil {
                    fmt.Println("Error: pathrider connect: "+args[1]+": "+err1.Error())
                }
                if err2!=nil {
                    fmt.Println("Error: pathrider connect: "+args[2]+": "+err2.Error())
                }
                if (err1==nil) && (err2==nil) {
                    fmt.Println("forwarding source nodes")
                    nodeSucc,edgeSucc=GetSuccessors(edges)
                    forward=ForwardEdges(sources,nodeSucc,edgeSucc)
                    fmt.Println("backwarding target nodes")
                    nodePred,edgePred=GetPredecessors(edges)
                    backward=BackwardEdges(targets,nodePred,edgePred)
                    if len(forward)==0 {
                        fmt.Println("Warning: pathrider connect: "+args[1]+": no forward paths found")
                    }
                    if len(backward)==0 {
                        fmt.Println("Warning: pathrider connect: "+args[2]+": no backward paths found")
                    }
                    if (len(forward)!=0) && (len(backward)!=0) {
                        fmt.Println("computing connecting paths")
                        intersect=IntersectEdges(forward,backward)
                        if len(intersect)==0 {
                            fmt.Println("Warning: pathrider connect: no connecting paths found")
                        } else {
                            fmt.Println("writing connecting paths: "+outFile)
                            err1=WriteNetwork(outFile,intersect,edgeNames)
                            if err1!=nil {
                                fmt.Println("Error: pathrider connect: "+outFile+": "+err1.Error())
                            } else if getShortest {
                                fmt.Println("computing shortest connecting paths")
                                noSelfLoop,selfLooped=RmSelfLoops(intersect)
                                nodeSucc,edgeSucc=GetSuccessors(noSelfLoop)
                                allShortest=AllShortestPaths(sources,targets,selfLooped,nodeSucc,edgeSucc)
                                outFilePath,outFileBase=filepath.Split(outFile)
                                outFileBase=strings.TrimSuffix(outFileBase,".sif")
                                outFileBase+="-shortest.sif"
                                outFile=filepath.Join(outFilePath,outFileBase)
                                fmt.Println("writing shortest connecting paths: "+outFile)
                                err1=WriteNetwork(outFile,allShortest,edgeNames)
                                if err1!=nil {
                                    fmt.Println("Error: pathrider connect: "+outFile+": "+err1.Error())
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
