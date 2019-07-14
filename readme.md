# Finding paths of interest in networks

Copyright 2019 [Arnaud Poret](https://github.com/arnaudporet)

This work is licensed under the [BSD 2-Clause License](https://raw.githubusercontent.com/arnaudporet/pathrider/master/license.txt).

## pathrider

pathrider is a tool for finding paths of interest in networks. It currently provides 2 commands:

* connect: find the paths connecting some nodes of interest in a network
* stream: find the upstream/downstream paths starting from some nodes of interest in a network

pathrider handles networks encoded in the SIF file format (see at the end of this readme file).

pathrider is implemented in [Go](https://golang.org) (see at the end of this readme file).

## Building/running

The simplest:

```sh
go run pathrider -help
```

A little bit more smarter:

```sh
go build pathrider.go
./pathrider -help
```

Even more smarter:

```sh
go build pathrider.go
mv pathrider /some/where/in/your/$PATH/
pathrider -help
```

Note that `go run` builds pathrider each time before running it, so building it is preferable. The built binaries can ultimately be placed somewhere in your `$PATH`.

The Go package can have different names depending on your operating system. For example, with [Ubuntu](https://ubuntu.com) the Go package is named `golang`. Consequently, running a Go file with Ubuntu might be `golang-go run yourfile.go` instead of `go run yourfile.go` with [Arch Linux](https://www.archlinux.org).

## Usage

### pathrider

Usage:

```
pathrider [options]
pathrider <command> [options] <arguments>
```

Options:

* non command-specific options:
    * `-l/-license`: print the BSD 2-Clause License under which pathrider is
    * `-h/-help`: print help
* command-specific options: see `pathrider <command> -help`

Commands:

* `connect`: find the paths connecting some nodes of interest in a network
* `stream`: find the upstream/downstream paths starting from some nodes of interest in a network

Arguments: see the command-specific help (`pathrider <command> -help`)

Cautions:

* pathrider handles networks encoded in the SIF file format (see at the end of this readme file)
* pathrider does not handle multi-edges (i.e. two or more edges having the same source and target nodes)
* note that duplicated edges are multi-edges
* edges are assumed to be directed

For command-specific help, run: `pathrider <command> -help`

### pathrider connect

Find the paths connecting some nodes of interest in a network.

Typical use is to find, in a network, the paths connecting some source nodes to some target nodes.

Usage:

```
pathrider connect [options] <networkFile> <sourceFile> <targetFile>
```

Positional arguments:

* `<networkFile>`: the network encoded in a SIF file
* `<sourceFile>`: the source nodes listed in a file (one node per line)
* `<targetFile>`: the target nodes listed in a file (one node per line)

Options:

* `-s/-shortest`: also find the shortest connecting paths (default: not used by default)
* `-o/-out <file>`: the output SIF file (default: `out.sif`)
* `-b/-blacklist <file>`: a file containing a list of nodes to be blacklisted (one node per line), the paths containing such nodes will not be considered (default: not used by default)
* `-h/-help`: print help

Output files (unless changed with `-o/-out`):

* `out.sif`: a SIF file encoding all the paths connecting the source nodes to the target nodes in the network
* `out-shortest.sif`: a SIF file encoding only the shortest connecting paths (requires `-s/-shortest`)

Cautions:

* the network must be in the SIF file format (see at the end of this readme file)
* the network must not contain multi-edges (i.e. two or more edges having the same source and target nodes)
* note that duplicated edges are multi-edges
* edges are assumed to be directed
* the source and target nodes must be listed in separate files with one node per line
* if sources = targets then provide the same node list twice

### pathrider stream

Find the upstream/downstream paths starting from some nodes of interest in a network.

Typical use is to find, in a network, the paths regulating some nodes (the upstream paths) or regulated by some nodes (the downstream paths).

Usage:

```
pathrider stream [options] <networkFile> <rootFile> <direction>
```

Positional arguments:

* `<networkFile>`: the network encoded in a SIF file
* `<rootFile>`: the root nodes listed in a file (one node per line)
* `<direction>`: follows the up stream (`up`) or the down stream (`down`)

Options:

* `-o/-out <file>`: the output SIF file (default: `out.sif`)
* `-b/-blacklist <file>`: a file containing a list of nodes to be blacklisted (one node per line), the paths containing such nodes will not be considered (default: not used by default)
* `-h/-help`: print help

Output file (unless changed with `-o/-out`):

* `out.sif`: a SIF file encoding the upstream/downstream paths starting from the root nodes in the network

Cautions:

* the network must be in the SIF file format (see at the end of this readme file)
* the network must not contain multi-edges (i.e. two or more edges having the same source and target nodes)
* note that duplicated edges are multi-edges
* edges are assumed to be directed
* the root nodes must be listed in a file with one node per line

## Examples

All the networks used in these examples are adapted from human signaling pathways coming from [KEGG Pathway](https://www.genome.jp/kegg/pathway.html).

### pathrider connect

* ErbB signaling pathway
    * `pathrider connect -s ErbB_signaling_pathway.sif sources.txt targets.txt`
    * networkFile: the ErbB signaling pathway (239 edges)
    * sourceFile: contains the nodes EGFR (i.e. ERBB1), ERBB2, ERBB3 and ERBB4
    * targetFile: contains the node MTOR
    * results:
        * out.sif (83 edges), also in SVG for visualization
        * out-shortest.sif (50 edges), also in SVG for visualization

* Insulin signaling pathway
    * `pathrider connect -s Insulin_signaling_pathway.sif sources.txt targets.txt`
    * networkFile: the insulin signaling pathway (407 edges)
    * sourceFile: contains the node INSR
    * targetFile: contains the nodes GSK3B and MAPK1
    * results:
        * out.sif (69 edges), also in SVG for visualization
        * out-shortest.sif (69 edges), also in SVG for visualization

* Cell cycle
    * `pathrider connect -s Cell_cycle.sif nodes.txt nodes.txt`
    * networkFile: the cell cycle (650 edges)
    * sourceFile: contains the node RB1
    * targetFile = sourceFile: for getting the paths connecting RB1 to itself
    * results:
        * out.sif (84 edges), also in SVG for visualization
        * out-shortest.sif (22 edges), also in SVG for visualization

* Cell survival
    * to illustrate the advantage of also computing the shortest connecting paths, this example is voluntarily bigger
    * it is made of the following human KEGG pathways: Apoptosis, Cell cycle, p53 signaling pathway, ErbB signaling pathway, TNF signaling pathway, TGF-beta signaling pathway, FoxO signaling pathway, Calcium signaling pathway, MAPK signaling pathway, PI3K-Akt signaling pathway and NF-kappa B signaling pathway
    * these pathways are involved in the cell growth/cell death balance
    * `pathrider connect -s Cell_survival.sif nodes.txt nodes.txt`
    * networkFile: some cell survival signaling pathways (11147 edges)
    * sourceFile: contains the nodes CASP3 (cell death effector), PIK3CA (involved in growth promoting signaling pathways) and TP53 (tumor suppressor)
    * targetFile = sourceFile: for viewing how these biological entities interact with each other
    * results:
        * out.sif (819 edges), also in SVG for a quite challenging visualization
        * out-shortest.sif (84 edges), also in SVG for an easier visualization, but only of the shortest connecting paths

### pathrider stream

* ErbB signaling pathway
    * `pathrider stream ErbB_signaling_pathway.sif roots.txt up`
    * networkFile: the ErbB signaling pathway (239 edges)
    * rootFile: contains the nodes JUN and MYC
    * direction: upstream
    * result: out.sif (133 edges), also in SVG for visualization

The ErbB signaling pathway is a growth-promoting signaling pathway typically activated by the epidermal growth factor (EGF).

JUN and MYC are two transcription factors influencing the expression of target genes following EGF stimulation.

The resulting file `out.sif` converted to SVG shows the upstream paths (i.e. the regulating paths) of JUN (red) and MYC (green) in the ErbB signaling pathway. It highlights that JUN and MYC share common elements in there regulating paths (red and green) and also specific elements (red or green). Note that other regulating paths outside of the ErbB signaling pathway exist.

* Toll-like receptor signaling pathway
    * `pathrider stream Toll-like_receptor_signaling_pathway.sif roots.txt down`
    * networkFile: the Toll-like receptor signaling pathway (219 edges)
    * rootFile: contains the nodes TLR3 and TLR4
    * direction: downstream
    * result: out.sif (152 edges), also in SVG for visualization

The Toll-like receptors (TLRs) are cell surface receptors which can be activated by various pathogen-associated molecular patterns (PAMPs).

PAMPs are molecules originating from microorganisms. TLRs are able to detect them in order to signal the presence of potentially harmful pathogens. There are several types of TLRs, each being able to detect specific PAMPs.

TLR3 can detect the presence of double-stranded RNA (dsRNA) coming from RNA viruses whereas TLR4 can detect lipopolysaccharide (LPS) coming from Gram-negative bacteria.

The resulting file `out.sif` converted to SVG shows the downstream paths (i.e. the effector paths) of TLR3 (red) and TLR4 (green). It highlights that TLR3 and TLR4 share common effectors (red and green) and also specific ones (red or green).

## The SIF file format

In a SIF file encoding a network, each line encodes an edge as follows:

```
source \t interaction \t target
```

Note that the field separator is the tabulation `\t`: the SIF file format is the tab-separated values format (TSV) with exactly 3 columns.

For example, the edge representing the activation of RAF1 by HRAS is a line of a SIF file encoded as follows:

```
HRAS \t activation \t RAF1
```

## Go

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:

* `go` (Arch Linux)
* `golang` (Ubuntu)

Otherwise, see https://golang.org/dl/ or https://golang.org/doc/install
