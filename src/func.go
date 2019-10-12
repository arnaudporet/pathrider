// Copyright (C) 2019 Arnaud Poret
// This work is licensed under the GNU General Public License.
// To view a copy of this license, visit https://www.gnu.org/licenses/gpl.html.

// WARNING The functions in the present file do not fully handle exceptions and
// errors. Instead, they assume that such handling is performed upstream by
// top-level functions of pathrider. Consequently, be careful if using them
// "as is" outside of pathrider.

package main
import (
    "encoding/csv"
    "errors"
    "os"
)
func AllShortestPaths(sources,targets,selfLooped []string,nodeSucc map[string][]string,edgeSucc map[string]map[string][][]string) [][]string {
    var (
        source,target string
        edge []string
        layers,shortest,allShortest [][]string
        nodePred map[string][]string
        edgePred map[string]map[string][][]string
    )
    for _,source=range sources {
        layers=GetLayers(source,nodeSucc,edgeSucc)
        nodePred,edgePred=GetPredecessors(layers)
        for _,target=range targets {
            if (source==target) && IsInList(selfLooped,source) {
                shortest=[][]string{[]string{source,target}}
            } else {
                shortest=ShortestPaths(source,target,nodePred,edgePred)
            }
            for _,edge=range shortest {
                if !IsInList2(allShortest,edge) {
                    allShortest=append(allShortest,CopyList(edge))
                }
            }
        }
    }
    return allShortest
}
func BackwardEdges(seeds []string,nodePred map[string][]string,edgePred map[string]map[string][][]string) [][]string {
    var (
        seed,npred string
        edge,epred []string
        backward,newCheck,toCheck [][]string
    )
    for _,seed=range seeds {
        for _,npred=range nodePred[seed] {
            backward=append(backward,[]string{npred,seed})
            newCheck=append(newCheck,[]string{npred,seed})
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
func ForwardEdges(seeds []string,nodeSucc map[string][]string,edgeSucc map[string]map[string][][]string) [][]string {
    var (
        seed,nsucc string
        edge,esucc []string
        forward,newCheck,toCheck [][]string
    )
    for _,seed=range seeds {
        for _,nsucc=range nodeSucc[seed] {
            forward=append(forward,[]string{seed,nsucc})
            newCheck=append(newCheck,[]string{seed,nsucc})
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
func GetLayers(seed string,nodeSucc map[string][]string,edgeSucc map[string]map[string][][]string) [][]string {
    var (
        nsucc string
        edge,esucc,visited []string
        layer,newLayer,edges [][]string
    )
    for _,nsucc=range nodeSucc[seed] {
        newLayer=append(newLayer,[]string{seed,nsucc})
        edges=append(edges,[]string{seed,nsucc})
    }
    for {
        for _,edge=range newLayer {
            visited=append(visited,edge[1])
        }
        layer=CopyList2(newLayer)
        newLayer=[][]string{}
        for _,edge=range layer {
            for _,esucc=range edgeSucc[edge[0]][edge[1]] {
                if !IsInList2(edges,esucc) && !IsInList(visited,esucc[1]) {
                    newLayer=append(newLayer,CopyList(esucc))
                    edges=append(edges,CopyList(esucc))
                }
            }
        }
        if len(newLayer)==0 {
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
func ReadNetwork(networkFile string) ([]string,[][]string,map[string]map[string][]string,error) {
    var (
        err error
        node string
        nodes,edge,line []string
        edges,lines [][]string
        edgeNames map[string]map[string][]string
        file *os.File
        reader *csv.Reader
    )
    edgeNames=make(map[string]map[string][]string)
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
                if !IsInList2(edges,edge) {
                    edges=append(edges,CopyList(edge))
                }
                for _,node=range edge {
                    if !IsInList(nodes,node) {
                        nodes=append(nodes,node)
                    }
                }
                edgeNames[line[0]]=make(map[string][]string)
            }
            if len(edges)==0 {
                err=errors.New("empty after reading")
            } else {
                for _,line=range lines {
                    edgeNames[line[0]][line[2]]=[]string{}
                }
                for _,line=range lines {
                    if !IsInList(edgeNames[line[0]][line[2]],line[1]) {
                        edgeNames[line[0]][line[2]]=append(edgeNames[line[0]][line[2]],line[1])
                    }
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
                    err=errors.New(line[0]+": node not in network")
                    break
                } else if !IsInList(nodes,line[0]) {
                    nodes=append(nodes,line[0])
                }
            }
            if err==nil {
                if len(nodes)==0 {
                    err=errors.New("empty after reading")
                }
            }
        }
    }
    return nodes,err
}
func RmNodes(edges [][]string,edgeNames map[string]map[string][]string,blackNodes []string) ([]string,[][]string,map[string]map[string][]string,error) {
    var (
        err error
        node string
        edge,newNodes []string
        newEdges [][]string
        newEdgeNames map[string]map[string][]string
    )
    newEdgeNames=make(map[string]map[string][]string)
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
            newEdgeNames[edge[0]]=make(map[string][]string)
        }
        for _,edge=range newEdges {
            newEdgeNames[edge[0]][edge[1]]=[]string{}
        }
        for _,edge=range newEdges {
            newEdgeNames[edge[0]][edge[1]]=CopyList(edgeNames[edge[0]][edge[1]])
        }
    }
    return newNodes,newEdges,newEdgeNames,err
}
func RmSelfLoops(edges [][]string) ([][]string,[]string) {
    var (
        edge,selfLooped []string
        noSelfLoop [][]string
    )
    for _,edge=range edges {
        if edge[0]==edge[1] {
            selfLooped=append(selfLooped,edge[0])
        } else {
            noSelfLoop=append(noSelfLoop,CopyList(edge))
        }
    }
    return noSelfLoop,selfLooped
}
func ShortestPaths(source,target string,nodePred map[string][]string,edgePred map[string]map[string][][]string) [][]string {
    var (
        found bool
        npred string
        edge,epred []string
        newCheck,toCheck,shortest [][]string
    )
    found=false
    for _,npred=range nodePred[target] {
        shortest=append(shortest,[]string{npred,target})
        newCheck=append(newCheck,[]string{npred,target})
    }
    for {
        for _,edge=range newCheck {
            if edge[0]==source {
                found=true
                break
            }
        }
        if found {
            break
        } else {
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
                shortest=[][]string{}
                break
            }
        }
    }
    return shortest
}
func TerminalNodes(nodeSP map[string][]string) []string {
    var (
        node string
        termNodes []string
    )
    for node=range nodeSP {
        if len(nodeSP[node])==0 {
            termNodes=append(termNodes,node)
        }
    }
    return termNodes
}
func WriteNetwork(networkFile string,edges [][]string,edgeNames map[string]map[string][]string) error {
    var (
        err error
        name string
        edge []string
        lines [][]string
        file *os.File
        writer *csv.Writer
    )
    for _,edge=range edges {
        for _,name=range edgeNames[edge[0]][edge[1]] {
            lines=append(lines,[]string{edge[0],name,edge[1]})
        }
    }
    if len(lines)==0 {
        err=errors.New("empty before writing")
    } else {
        file,err=os.Create(networkFile)
        defer file.Close()
        if err==nil {
            writer=csv.NewWriter(file)
            writer.Comma='\t'
            writer.UseCRLF=false
            err=writer.WriteAll(lines)
        }
    }
    return err
}
