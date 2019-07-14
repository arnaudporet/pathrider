// Copyright 2019 Arnaud Poret
// This work is licensed under the BSD 2-Clause License.
package main
import (
    "encoding/csv"
    "errors"
    "flag"
    "fmt"
    "os"
    "path/filepath"
    "strings"
)
//############################################################################//
//#### MAIN ##################################################################//
//############################################################################//
func main() {
    var (
        err error
        help,license bool
        args []string
        flagSet *flag.FlagSet
    )
    flagSet=flag.NewFlagSet("",flag.ContinueOnError)
    flagSet.Usage=func() {}
    flagSet.BoolVar(&help,"help",false,"")
    flagSet.BoolVar(&help,"h",false,"")
    flagSet.BoolVar(&license,"license",false,"")
    flagSet.BoolVar(&license,"l",false,"")
    err=flagSet.Parse(os.Args[1:])
    if err!=nil {
        fmt.Println("Error: pathrider: "+err.Error())
    } else if help {
        fmt.Println(strings.Join([]string{
            "",
            "pathrider is a tool for finding paths of interest in networks.",
            "",
            "pathrider currently provides 2 commands:",
            "    * connect: find the paths connecting some nodes of interest in a network",
            "    * stream: find the upstream/downstream paths starting from some nodes of interest in a network",
            "",
            "For command-specific help, run: pathrider <command> -help",
            "",
            "Cautions:",
            "    * pathrider handles networks encoded in the SIF file format (see the readme file of pathrider)",
            "    * pathrider does not handle multi-edges (i.e. two or more edges having the same source and target nodes)",
            "    * note that duplicated edges are multi-edges",
            "    * edges are assumed to be directed",
            "",
            "Usage:",
            "    * pathrider [options]",
            "    * pathrider <command> [options] <arguments>",
            "",
            "Options:",
            "    * non command-specific options:",
            "        * -l/-license: print the BSD 2-Clause License under which pathrider is",
            "        * -h/-help: print this help",
            "    * command-specific options: pathrider <command> -help",
            "",
            "Commands:",
            "    * connect: find the paths connecting some nodes of interest in a network",
            "    * stream: find the upstream/downstream paths starting from some nodes of interest in a network",
            "",
            "Arguments: see the command-specific help (pathrider <command> -help)",
            "",
            "For more information, see https://github.com/arnaudporet/pathrider",
            "",
        },"\n"))
    } else if license {
        fmt.Println(strings.Join([]string{
            "",
            "Copyright 2019 Arnaud Poret",
            "",
            "Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:",
            "",
            "1. Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.",
            "",
            "2. Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.",
            "",
            "THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS \"AS IS\" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.",
            "",
        },"\n"))
    } else if len(flagSet.Args())==0 {
        fmt.Println("Error: pathrider: missing command, expecting one of: connect, stream")
    } else {
        args=flagSet.Args()
        if args[0]=="connect" {
            Connect()
        } else if args[0]=="stream" {
            Stream()
        } else {
            fmt.Println("Error: pathrider: "+args[0]+": unknown command, expecting one of: connect, stream")
        }
    }
}
//############################################################################//
//#### COMMAND ###############################################################//
//############################################################################//
func Connect() {
    var (
        err1,err2,err3 error
        help,getShortest bool
        outFile,outFilePath,outFileBase,shortestFile,blackFile string
        args,nodes,sources,targets,blackNodes []string
        edges,forward,backward,intersect,allShortest [][]string
        edgeNames map[string]map[string]string
        flagSet *flag.FlagSet
    )
    flagSet=flag.NewFlagSet("",flag.ContinueOnError)
    flagSet.Usage=func() {}
    flagSet.BoolVar(&help,"help",false,"")
    flagSet.BoolVar(&help,"h",false,"")
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
            "Typical use is to find, in a network, the paths connecting some source nodes to some target nodes.",
            "",
            "Cautions:",
            "    * the network must be in the SIF file format (see the readme file of pathrider)",
            "    * the network must not contain multi-edges (i.e. two or more edges having the same source and target nodes)",
            "    * note that duplicated edges are multi-edges",
            "    * edges are assumed to be directed",
            "    * the source and target nodes must be listed in separate files with one node per line",
            "    * if sources = targets then provide the same node list twice",
            "",
            "Output files (unless changed with -o/-out):",
            "    * out.sif: a SIF file encoding all the paths connecting the source nodes to the target nodes in the network",
            "    * out-shortest.sif: a SIF file encoding only the shortest connecting paths (requires -s/-shortest)",
            "",
            "Usage: pathrider connect [options] <networkFile> <sourceFile> <targetFile>",
            "",
            "Positional arguments:",
            "    * <networkFile>: the network encoded in a SIF file",
            "    * <sourceFile>: the source nodes listed in a file (one node per line)",
            "    * <targetFile>: the target nodes listed in a file (one node per line)",
            "",
            "Options:",
            "    * -s/-shortest: also find the shortest connecting paths (default: not used by default)",
            "    * -o/-out <file>: the output SIF file (default: out.sif)",
            "    * -b/-blacklist <file>: a file containing a list of nodes to be blacklisted (one node per line), the paths containing such nodes will not be considered (default: not used by default)",
            "    * -h/-help: print this help",
            "",
            "For more information, see https://github.com/arnaudporet/pathrider",
            "",
        },"\n"))
    } else if filepath.Ext(outFile)!=".sif" {
        fmt.Println("Error: pathrider connect: "+outFile+": must have the \".sif\" file extension")
    } else if len(flagSet.Args())!=3 {
        fmt.Println("Error: pathrider connect: wrong number of positional arguments, expecting: <networkFile> <sourceFile> <targetFile>")
    } else {
        args=flagSet.Args()
        fmt.Println("reading "+args[0])
        nodes,edges,edgeNames,err1=ReadNetwork(args[0])
        if err1!=nil {
            fmt.Println("Error: pathrider connect: "+args[0]+": "+err1.Error())
        } else {
            fmt.Println("reading "+args[1])
            sources,err1=ReadNodes(args[1],nodes)
            fmt.Println("reading "+args[2])
            targets,err2=ReadNodes(args[2],nodes)
            if blackFile!="" {
                fmt.Println("reading "+blackFile)
                blackNodes,err3=ReadNodes(blackFile,nodes)
            }
            if err1!=nil {
                fmt.Println("Error: pathrider connect: "+args[1]+": "+err1.Error())
            }
            if err2!=nil {
                fmt.Println("Error: pathrider connect: "+args[2]+": "+err2.Error())
            }
            if err3!=nil {
                fmt.Println("Error: pathrider connect: "+blackFile+": "+err3.Error())
            }
            if (err1==nil) && (err2==nil) && (err3==nil) {
                if blackFile!="" {
                    fmt.Println("blacklisting from "+blackFile)
                    nodes,edges,edgeNames,err1=RmNodes(edges,edgeNames,blackNodes)
                }
                if err1!=nil {
                    fmt.Println("Error: pathrider connect: "+blackFile+": "+err1.Error())
                } else {
                    fmt.Println("forwarding "+args[1])
                    forward=ForwardEdges(sources,edges)
                    fmt.Println("backwarding "+args[2])
                    backward=BackwardEdges(targets,edges)
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
                            fmt.Println("Warning: pathrider connect: "+args[1]+" "+args[2]+": no connecting paths found")
                        } else {
                            fmt.Println("writing "+outFile)
                            err1=WriteNetwork(outFile,intersect,edgeNames)
                            if err1!=nil {
                                fmt.Println("Error: pathrider connect: "+outFile+": "+err1.Error())
                            } else if getShortest {
                                fmt.Println("computing shortest connecting paths")
                                allShortest=AllShortestPaths(sources,targets,intersect)
                                if len(allShortest)==0 {
                                    fmt.Println("Warning: pathrider connect: "+args[1]+" "+args[2]+": no shortest connecting paths found")
                                } else {
                                    outFilePath,outFileBase=filepath.Split(outFile)
                                    outFileBase=strings.TrimSuffix(outFileBase,".sif")
                                    outFileBase+="-shortest.sif"
                                    shortestFile=filepath.Join(outFilePath,outFileBase)
                                    fmt.Println("writing "+shortestFile)
                                    err1=WriteNetwork(shortestFile,allShortest,edgeNames)
                                    if err1!=nil {
                                        fmt.Println("Error: pathrider connect: "+shortestFile+": "+err1.Error())
                                    }
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}
func Stream() {
    var (
        err1,err2 error
        help bool
        outFile,blackFile string
        args,nodes,roots,blackNodes []string
        edges,ward [][]string
        edgeNames map[string]map[string]string
        flagSet *flag.FlagSet
    )
    flagSet=flag.NewFlagSet("",flag.ContinueOnError)
    flagSet.Usage=func() {}
    flagSet.BoolVar(&help,"help",false,"")
    flagSet.BoolVar(&help,"h",false,"")
    flagSet.StringVar(&outFile,"out","out.sif","")
    flagSet.StringVar(&outFile,"o","out.sif","")
    flagSet.StringVar(&blackFile,"blacklist","","")
    flagSet.StringVar(&blackFile,"b","","")
    err1=flagSet.Parse(os.Args[2:])
    if err1!=nil {
        fmt.Println("Error: pathrider stream: "+err1.Error())
    } else if help {
        fmt.Println(strings.Join([]string{
            "",
            "Find the upstream/downstream paths starting from some nodes of interest in a network.",
            "",
            "Typical use is to find, in a network, the paths regulating some nodes (the upstream paths) or regulated by some nodes (the downstream paths).",
            "",
            "Cautions:",
            "    * the network must be in the SIF file format (see the readme file of pathrider)",
            "    * the network must not contain multi-edges (i.e. two or more edges having the same source and target nodes)",
            "    * note that duplicated edges are multi-edges",
            "    * edges are assumed to be directed",
            "    * the root nodes must be listed in a file with one node per line",
            "",
            "Output file (unless changed with -o/-out):",
            "    * out.sif: a SIF file encoding the upstream/downstream paths starting from the root nodes in the network",
            "",
            "Usage: pathrider stream [options] <networkFile> <rootFile> <direction>",
            "",
            "Positional arguments:",
            "    * <networkFile>: the network encoded in a SIF file",
            "    * <rootFile>: the root nodes listed in a file (one node per line)",
            "    * <direction>: follows the up stream (up) or the down stream (down)",
            "",
            "Options:",
            "    * -o/-out <file>: the output SIF file (default: out.sif)",
            "    * -b/-blacklist <file>: a file containing a list of nodes to be blacklisted (one node per line), the paths containing such nodes will not be considered (default: not used by default)",
            "    * -h/-help: print this help",
            "",
            "For more information, see https://github.com/arnaudporet/pathrider",
            "",
        },"\n"))
    } else if filepath.Ext(outFile)!=".sif" {
        fmt.Println("Error: pathrider stream: "+outFile+": must have the \".sif\" file extension")
    } else if len(flagSet.Args())!=3 {
        fmt.Println("Error: pathrider stream: wrong number of positional arguments, expecting: <networkFile> <rootFile> <direction>")
    } else if !IsInList([]string{"up","down"},flagSet.Arg(2)) {
        fmt.Println("Error: pathrider stream: "+flagSet.Arg(2)+": unknown direction, expecting one of: up, down")
    } else {
        args=flagSet.Args()
        fmt.Println("reading "+args[0])
        nodes,edges,edgeNames,err1=ReadNetwork(args[0])
        if err1!=nil {
            fmt.Println("Error: pathrider stream: "+args[0]+": "+err1.Error())
        } else {
            fmt.Println("reading "+args[1])
            roots,err1=ReadNodes(args[1],nodes)
            if blackFile!="" {
                fmt.Println("reading "+blackFile)
                blackNodes,err2=ReadNodes(blackFile,nodes)
            }
            if err1!=nil {
                fmt.Println("Error: pathrider stream: "+args[1]+": "+err1.Error())
            }
            if err2!=nil {
                fmt.Println("Error: pathrider stream: "+blackFile+": "+err2.Error())
            }
            if (err1==nil) && (err2==nil) {
                if blackFile!="" {
                    fmt.Println("blacklisting from "+blackFile)
                    nodes,edges,edgeNames,err1=RmNodes(edges,edgeNames,blackNodes)
                }
                if err1!=nil {
                    fmt.Println("Error: pathrider stream: "+blackFile+": "+err1.Error())
                } else {
                    fmt.Println(args[2]+"streaming "+args[1])
                    if args[2]=="down" {
                        ward=ForwardEdges(roots,edges)
                    } else if args[2]=="up" {
                        ward=BackwardEdges(roots,edges)
                    }
                    if len(ward)==0 {
                        fmt.Println("Warning: pathrider stream: "+args[1]+": no "+args[2]+"stream paths found")
                    } else {
                        fmt.Println("writing "+outFile)
                        err1=WriteNetwork(outFile,ward,edgeNames)
                        if err1!=nil {
                            fmt.Println("Error: pathrider stream: "+outFile+": "+err1.Error())
                        }
                    }
                }
            }
        }
    }
}
//############################################################################//
//#### FUNCTION ##############################################################//
//############################################################################//
func AllShortestPaths(sources,targets []string,edges [][]string) [][]string {
    var (
        source,target string
        selfLooped,edge []string
        noSelfLoops,layers,shortest,allShortest [][]string
        nodeSucc,nodePred map[string][]string
        edgeSucc,edgePred map[string]map[string][][]string
    )
    noSelfLoops,selfLooped=RmSelfLoops(edges)
    nodeSucc,edgeSucc=GetSuccessors(noSelfLoops)
    for _,source=range sources {
        layers=GetLayers(source,nodeSucc,edgeSucc)
        nodePred,edgePred=GetPredecessors(layers)
        for _,target=range targets {
            shortest=ShortestPaths(source,target,selfLooped,nodePred,edgePred)
            for _,edge=range shortest {
                if !IsInList2(allShortest,edge) {
                    allShortest=append(allShortest,CopyList(edge))
                }
            }
        }
    }
    return allShortest
}
func BackwardEdges(roots []string,edges [][]string) [][]string {
    var (
        root,npred string
        edge,epred []string
        toCheck,newCheck,backward [][]string
        nodePred map[string][]string
        edgePred map[string]map[string][][]string
    )
    nodePred,edgePred=GetPredecessors(edges)
    for _,root=range roots {
        for _,npred=range nodePred[root] {
            backward=append(backward,[]string{npred,root})
            newCheck=append(newCheck,[]string{npred,root})
        }
    }
    for {
        toCheck=CopyList2(newCheck)
        newCheck=[][]string{}
        for _,edge=range toCheck {
            for _,epred=range edgePred[edge[0]][edge[1]] {
                if !IsInList2(backward,epred) {
                    backward=append(backward,CopyList(epred))
                    newCheck=append(newCheck,CopyList(epred))
                }
            }
        }
        if len(newCheck)==0 {
            break
        }
    }
    return backward
}
func CopyList(list []string) []string {
    var y []string
    y=make([]string,len(list))
    copy(y,list)
    return y
}
func CopyList2(list2 [][]string) [][]string {
    var (
        i int
        y [][]string
    )
    y=make([][]string,len(list2))
    for i=range list2 {
        y[i]=make([]string,len(list2[i]))
        copy(y[i],list2[i])
    }
    return y
}
func ForwardEdges(roots []string,edges [][]string) [][]string {
    var (
        root,nsucc string
        edge,esucc []string
        toCheck,newCheck,forward [][]string
        nodeSucc map[string][]string
        edgeSucc map[string]map[string][][]string
    )
    nodeSucc,edgeSucc=GetSuccessors(edges)
    for _,root=range roots {
        for _,nsucc=range nodeSucc[root] {
            forward=append(forward,[]string{root,nsucc})
            newCheck=append(newCheck,[]string{root,nsucc})
        }
    }
    for {
        toCheck=CopyList2(newCheck)
        newCheck=[][]string{}
        for _,edge=range toCheck {
            for _,esucc=range edgeSucc[edge[0]][edge[1]] {
                if !IsInList2(forward,esucc) {
                    forward=append(forward,CopyList(esucc))
                    newCheck=append(newCheck,CopyList(esucc))
                }
            }
        }
        if len(newCheck)==0 {
            break
        }
    }
    return forward
}
func GetLayers(root string,nodeSucc map[string][]string,edgeSucc map[string]map[string][][]string) [][]string {
    var (
        nsucc string
        edge,esucc,visited []string
        layer,edges [][]string
        layers [][][]string
    )
    for _,nsucc=range nodeSucc[root] {
        layer=append(layer,[]string{root,nsucc})
        edges=append(edges,[]string{root,nsucc})
    }
    for {
        layers=append(layers,CopyList2(layer))
        for _,edge=range layer {
            visited=append(visited,edge[1])
        }
        layer=[][]string{}
        for _,edge=range layers[len(layers)-1] {
            for _,esucc=range edgeSucc[edge[0]][edge[1]] {
                if !IsInList2(edges,esucc) && !IsInList(visited,esucc[1]) {
                    layer=append(layer,CopyList(esucc))
                    edges=append(edges,CopyList(esucc))
                }
            }
        }
        if len(layer)==0 {
            break
        }
    }
    return edges
}
func GetPredecessors(edges [][]string) (map[string][]string,map[string]map[string][][]string) {
    var (
        node,node2,node3 string
        edge []string
        nodePred map[string][]string
        edgePred map[string]map[string][][]string
    )
    nodePred=make(map[string][]string)
    edgePred=make(map[string]map[string][][]string)
    for _,edge=range edges {
        for _,node=range edge {
            nodePred[node]=[]string{}
        }
        edgePred[edge[0]]=make(map[string][][]string)
    }
    for _,edge=range edges {
        nodePred[edge[1]]=append(nodePred[edge[1]],edge[0])
        edgePred[edge[0]][edge[1]]=[][]string{}
    }
    for node=range nodePred {
        for _,node2=range nodePred[node] {
            for _,node3=range nodePred[node2] {
                edgePred[node2][node]=append(edgePred[node2][node],[]string{node3,node2})
            }
        }
    }
    return nodePred,edgePred
}
func GetSuccessors(edges [][]string) (map[string][]string,map[string]map[string][][]string) {
    var (
        node,node2,node3 string
        edge []string
        nodeSucc map[string][]string
        edgeSucc map[string]map[string][][]string
    )
    nodeSucc=make(map[string][]string)
    edgeSucc=make(map[string]map[string][][]string)
    for _,edge=range edges {
        for _,node=range edge {
            nodeSucc[node]=[]string{}
        }
        edgeSucc[edge[0]]=make(map[string][][]string)
    }
    for _,edge=range edges {
        nodeSucc[edge[0]]=append(nodeSucc[edge[0]],edge[1])
        edgeSucc[edge[0]][edge[1]]=[][]string{}
    }
    for node=range nodeSucc {
        for _,node2=range nodeSucc[node] {
            for _,node3=range nodeSucc[node2] {
                edgeSucc[node][node2]=append(edgeSucc[node][node2],[]string{node2,node3})
            }
        }
    }
    return nodeSucc,edgeSucc
}
func IntersectEdges(edges1,edges2 [][]string) [][]string {
    var (
        edge []string
        intersect [][]string
    )
    for _,edge=range edges1 {
        if IsInList2(edges2,edge) {
            intersect=append(intersect,CopyList(edge))
        }
    }
    return intersect
}
func IsInList(list []string,thatElement string) bool {
    var element string
    for _,element=range list {
        if element==thatElement {
            return true
        }
    }
    return false
}
func IsInList2(list2 [][]string,thatList []string) bool {
    var (
        found bool
        i int
        list []string
    )
    for _,list=range list2 {
        if len(list)==len(thatList) {
            found=true
            for i=range list {
                if list[i]!=thatList[i] {
                    found=false
                    break
                }
            }
            if found {
                return true
            }
        }
    }
    return false
}
func ReadNetwork(networkFile string) ([]string,[][]string,map[string]map[string]string,error) {
    var (
        err error
        node string
        nodes,edge,line []string
        edges,lines [][]string
        edgeNames map[string]map[string]string
        file *os.File
        reader *csv.Reader
    )
    edgeNames=make(map[string]map[string]string)
    file,err=os.Open(networkFile)
    defer file.Close()
    if err==nil {
        reader=csv.NewReader(file)
        reader.Comma='\t'
        reader.Comment=0
        reader.FieldsPerRecord=3
        reader.LazyQuotes=false
        reader.TrimLeadingSpace=true
        reader.ReuseRecord=true
        lines,err=reader.ReadAll()
        if err==nil {
            for _,line=range lines {
                edge=[]string{line[0],line[2]}
                if IsInList2(edges,edge) {
                    err=errors.New("multi-edges (or duplicated edges)")
                    break
                } else {
                    edges=append(edges,CopyList(edge))
                    for _,node=range edge {
                        if !IsInList(nodes,node) {
                            nodes=append(nodes,node)
                        }
                    }
                    edgeNames[line[0]]=make(map[string]string)
                }
            }
            if err==nil {
                for _,line=range lines {
                    edgeNames[line[0]][line[2]]=line[1]
                }
                if len(edges)==0 {
                    err=errors.New("empty after reading")
                }
            }
        }
    }
    return nodes,edges,edgeNames,err
}
func ReadNodes(nodeFile string,networkNodes []string) ([]string,error) {
    var (
        err error
        line,nodes []string
        lines [][]string
        file *os.File
        reader *csv.Reader
    )
    file,err=os.Open(nodeFile)
    defer file.Close()
    if err==nil {
        reader=csv.NewReader(file)
        reader.Comma='\t'
        reader.Comment=0
        reader.FieldsPerRecord=1
        reader.LazyQuotes=false
        reader.TrimLeadingSpace=true
        reader.ReuseRecord=true
        lines,err=reader.ReadAll()
        if err==nil {
            for _,line=range lines {
                if !IsInList(networkNodes,line[0]) {
                    fmt.Println("Warning: pathrider: "+nodeFile+": "+line[0]+" not in network")
                } else if !IsInList(nodes,line[0]) {
                    nodes=append(nodes,line[0])
                }
            }
            if len(nodes)==0 {
                err=errors.New("empty after reading")
            }
        }
    }
    return nodes,err
}
func RmNodes(edges [][]string,edgeNames map[string]map[string]string,blackNodes []string) ([]string,[][]string,map[string]map[string]string,error) {
    var (
        err error
        node string
        edge,newNodes []string
        newEdges [][]string
        newEdgeNames map[string]map[string]string
    )
    newEdgeNames=make(map[string]map[string]string)
    for _,edge=range edges {
        if !IsInList(blackNodes,edge[0]) && !IsInList(blackNodes,edge[1]) {
            newEdges=append(newEdges,CopyList(edge))
        }
    }
    if len(newEdges)==0 {
        err=errors.New("network empty after blacklisting")
    } else {
        for _,edge=range newEdges {
            for _,node=range edge {
                if !IsInList(newNodes,node) {
                    newNodes=append(newNodes,node)
                }
            }
            newEdgeNames[edge[0]]=make(map[string]string)
        }
        for _,edge=range newEdges {
            newEdgeNames[edge[0]][edge[1]]=edgeNames[edge[0]][edge[1]]
        }
    }
    return newNodes,newEdges,newEdgeNames,err
}
func RmSelfLoops(edges [][]string) ([][]string,[]string) {
    var (
        edge,selfLooped []string
        noSelfLoops [][]string
    )
    for _,edge=range edges {
        if edge[0]==edge[1] {
            selfLooped=append(selfLooped,edge[0])
        } else {
            noSelfLoops=append(noSelfLoops,CopyList(edge))
        }
    }
    return noSelfLoops,selfLooped
}
func ShortestPaths(source,target string,selfLooped []string,nodePred map[string][]string,edgePred map[string]map[string][][]string) [][]string {
    var (
        npred string
        edge,epred []string
        newCheck,toCheck,shortest [][]string
    )
    if (source==target) && IsInList(selfLooped,source) {
        return [][]string{[]string{source,target}}
    } else {
        for _,npred=range nodePred[target] {
            shortest=append(shortest,[]string{npred,target})
            newCheck=append(newCheck,[]string{npred,target})
        }
        for {
            for _,edge=range newCheck {
                if edge[0]==source {
                    return shortest
                }
            }
            toCheck=CopyList2(newCheck)
            newCheck=[][]string{}
            for _,edge=range toCheck {
                for _,epred=range edgePred[edge[0]][edge[1]] {
                    if !IsInList2(shortest,epred) {
                        shortest=append(shortest,CopyList(epred))
                        newCheck=append(newCheck,CopyList(epred))
                    }
                }
            }
            if len(newCheck)==0 {
                break
            }
        }
        return [][]string{}
    }
}
func WriteNetwork(networkFile string,edges [][]string,edgeNames map[string]map[string]string) error {
    var (
        err error
        edge []string
        lines [][]string
        file *os.File
        writer *csv.Writer
    )
    file,err=os.Create(networkFile)
    defer file.Close()
    if err==nil {
        for _,edge=range edges {
            lines=append(lines,[]string{edge[0],edgeNames[edge[0]][edge[1]],edge[1]})
        }
        writer=csv.NewWriter(file)
        writer.Comma='\t'
        writer.UseCRLF=false
        err=writer.WriteAll(lines)
    }
    return err
}
