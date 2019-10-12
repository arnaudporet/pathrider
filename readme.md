# Finding paths of interest in networks

Copyright 2019 [Arnaud Poret](https://github.com/arnaudporet)

This work is licensed under the [GNU General Public License](https://www.gnu.org/licenses/gpl-3.0.html).

## pathrider

pathrider is a tool for finding paths of interest in networks. It currently provides 2 commands:

* `connect`: find the paths connecting some nodes of interest in a network
* `stream`: find the upstream/downstream paths starting from some nodes of interest in a network

pathrider handles networks encoded in the SIF file format (see at the end of this readme file).

pathrider is implemented in [Go](https://golang.org) (see at the end of this readme file).

## Building and running

The Go package can have different names depending on the used operating system. For example, with [Ubuntu](https://ubuntu.com) the Go package is named `golang`. Consequently, calling the Go compiler might be `golang-go` instead of `go` with [Arch Linux](https://www.archlinux.org).

```sh
cd pathrider/ # move to the pathrider directory
go build -o pathrider ./src/ # build pathrider using its source directory with the Go compiler
./pathrider -help # run pathrider
```

Note that the path to the source directory `src/` supplied to the GO compiler `go` must be a relative path (_i.e._ `./src/`, not `src/` or `/absolute/path/to/src/`).

Once pathrider built, it can be moved somewhere into the `$PATH` to make it easily callable from everywhere:

```sh
mv pathrider /somewhere/in/the/$PATH/
pathrider -help
```

## Usage

### pathrider

Usage:

```
pathrider [options]
pathrider <command> [options] <arguments>
```

Positional argument:

* `<command>`: `connect`, `stream`

Options:

* `-l/-license`: print the GNU General Public License under which pathrider is
* `-u/-usage`: print usage only
* `-h/-help`: print help

For command-specific help, run `pathrider <command> -help`.

### pathrider connect

Find the paths connecting some nodes of interest in a network.

Typical use is to find in a network the paths connecting some source nodes to some target nodes.

Usage:

```
pathrider connect [options] <networkFile> <sourceFile> <targetFile>
```

Positional arguments:

* `<networkFile>`: the network encoded in a SIF file
* `<sourceFile>`: the source nodes listed in a file (one node per line)
* `<targetFile>`: the target nodes listed in a file (one node per line)
* if sources = targets then provide the same node list twice

Options:

* `-s/-shortest`: also find the shortest connecting paths (default: not used by default)
* `-b/-blacklist <file>`: a file containing a list of nodes to be blacklisted (one node per line), the paths containing such nodes will not be considered (default: not used by default)
* `-o/-out <file>`: the output SIF file (default: `out.sif`)
* `-u/-usage`: print usage only
* `-h/-help`: print help

Output file(s) (unless changed with `-o/-out`):

* `out.sif`: a SIF file encoding all the paths connecting the source nodes to the target nodes in the network
* `out-shortest.sif`: a SIF file encoding only the shortest connecting paths (requires `-s/-shortest`)

Cautions:

* the network must be in the SIF file format (see at the end of this readme file)
* edge duplicates are automatically removed
* edges are assumed to be directed

### pathrider stream

Find the upstream/downstream paths starting from some nodes of interest (the seed nodes) in a network.

Typical use is to find in a network the paths regulating some nodes (the upstream paths) or the paths regulated by some nodes (the downstream paths).

Usage:

```
pathrider stream [options] <networkFile> <seedFile> <direction>
```

Positional arguments:

* `<networkFile>`: the network encoded in a SIF file
* `<seedFile>`: the seed nodes listed in a file (one node per line)
* `<direction>`: follow the up stream (`up`) or the down stream (`down`)

Options:

* `-t/-terminal`: also find the terminal nodes reachable from the seed nodes, namely the nodes having no predecessors in case of upstreaming, or the nodes having no successors in case of downstreaming (default: not used by default)
* `-b/-blacklist <file>`: a file containing a list of nodes to be blacklisted (one node per line), the paths containing such nodes will not be considered (default: not used by default)
* `-o/-out <file>`: the output SIF file (default: `out.sif`)
* `-u/-usage`: print usage only
* `-h/-help`: print help

Output file(s) (unless changed with `-o/-out`):

* `out.sif`: a SIF file encoding the upstream/downstream paths starting from the seed nodes in the network
* `out-terminal.txt`: a file listing the upstream/downstream terminal nodes reachable from the seed nodes in the network (requires `-t/-terminal`)

Cautions:

* the network must be in the SIF file format (see at the end of this readme file)
* edge duplicates are automatically removed
* edges are assumed to be directed

## Examples

All the networks used in these examples are adapted from human signaling pathways coming from [KEGG Pathway](https://www.genome.jp/kegg/pathway.html) using [kgml2sif](https://github.com/arnaudporet/kgml2sif).

### pathrider connect

#### ErbB signaling pathway

* `pathrider connect -s ErbB_signaling_pathway.sif sources.txt targets.txt`
* networkFile: the ErbB signaling pathway (239 edges)
* sourceFile: contains the nodes EGFR (_i.e._ ERBB1), ERBB2, ERBB3 and ERBB4
* targetFile: contains the node MTOR
* results:
    * out.sif (83 edges), also in SVG for visualization
    * out-shortest.sif (50 edges), also in SVG for visualization

#### Insulin signaling pathway

* `pathrider connect -s Insulin_signaling_pathway.sif sources.txt targets.txt`
* networkFile: the insulin signaling pathway (429 edges)
* sourceFile: contains the node INSR
* targetFile: contains the nodes GSK3B and MAPK1
* results:
    * out.sif (69 edges), also in SVG for visualization
    * out-shortest.sif (69 edges), also in SVG for visualization

#### Cell cycle

* `pathrider connect -s Cell_cycle.sif nodes.txt nodes.txt`
* networkFile: the cell cycle (650 edges)
* sourceFile: contains the node RB1
* targetFile = sourceFile: for getting the paths connecting RB1 to itself
* results:
    * out.sif (84 edges), also in SVG for visualization
    * out-shortest.sif (22 edges), also in SVG for visualization

#### Cell survival

* to illustrate the interest of also computing the shortest connecting paths, this example is voluntarily bigger
* it is made of the following KEGG human signaling pathways involved in the cell growth/cell death balance:
    * Cell cycle
    * Apoptosis
    * p53 signaling pathway
    * TNF signaling pathway
    * TGF-beta signaling pathway
    * ErbB signaling pathway
    * MAPK signaling pathway
    * Calcium signaling pathway
    * PI3K-Akt signaling pathway
    * NF-kappa B signaling pathway
    * mTOR signaling pathway
    * FoxO signaling pathway
    * Phosphatidylinositol signaling system
* `pathrider connect -s Cell_survival.sif nodes.txt nodes.txt`
* networkFile: some cell survival signaling pathways (14 181 edges)
* sourceFile: contains the nodes CASP3 (cell death effector), PIK3CA (involved in growth promoting signaling pathways) and TP53 (tumor suppressor)
* targetFile = sourceFile: for viewing how these biological entities interact with each other
* results:
    * out.sif (2 870 edges), also in SVG for a quite challenging visualization
    * out-shortest.sif (106 edges), also in SVG for an easier visualization, but only of the shortest connecting paths

### pathrider stream

#### ErbB signaling pathway

* `pathrider stream ErbB_signaling_pathway.sif seeds.txt up`
* networkFile: the ErbB signaling pathway (239 edges)
* seedFile: contains the nodes JUN and MYC
* direction: upstream
* result: out.sif (133 edges), also in SVG for visualization

The ErbB signaling pathway is a growth-promoting signaling pathway typically activated by the epidermal growth factor (EGF).

JUN and MYC are two transcription factors influencing the expression of target genes following EGF stimulation.

The resulting file `out.sif` converted to SVG shows the upstream paths (_i.e._ the regulating paths) of JUN (green) and MYC (red) in the ErbB signaling pathway. It highlights that JUN and MYC share common elements in there regulating paths (red and green) and also specific elements (red or green). Note that other regulating paths outside of the ErbB signaling pathway exist.

#### Toll-like receptor signaling pathway

* `pathrider stream Toll-like_receptor_signaling_pathway.sif seeds.txt down`
* networkFile: the Toll-like receptor signaling pathway (222 edges)
* seedFile: contains the nodes TLR3 and TLR4
* direction: downstream
* result: out.sif (155 edges), also in SVG for visualization

The Toll-like receptors (TLRs) are cell surface receptors which can be activated by various pathogen-associated molecular patterns (PAMPs).

PAMPs are molecules originating from microorganisms. TLRs are able to detect them in order to signal the presence of potentially harmful pathogens. There are several types of TLRs, each being able to detect specific PAMPs.

TLR3 can detect the presence of double-stranded RNA (dsRNA) coming from RNA viruses whereas TLR4 can detect lipopolysaccharide (LPS) coming from Gram-negative bacteria.

The resulting file `out.sif` converted to SVG shows the downstream paths (_i.e._ the effector paths) of TLR3 (green) and TLR4 (red). It highlights that TLR3 and TLR4 share common effectors (red and green) and also specific ones (red or green).

## The SIF file format

In a SIF file encoding a network, each line encodes an edge as follows:

```
source \t interaction \t target
```

Note that the field separator is the tabulation: the SIF file format is the tab-separated values format (TSV) with exactly 3 columns.

For example, the edge representing the activation of RAF1 by HRAS is a line of a SIF file encoded as follows:

```
HRAS \t activation \t RAF1
```

## Go

Most [Linux distributions](https://distrowatch.com) provide Go in their official repositories. For example:

* `go` (Arch Linux)
* `golang` (Ubuntu)

Otherwise, see https://golang.org/doc/install or https://golang.org/dl/
