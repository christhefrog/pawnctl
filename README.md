##
<div align="center">
<img alt="pawnctl" src="logo.png" width=535/>

`pawnctl` is a minimalistic cli tool for automating pawn script compilation processes \
<sup>_pawnctl is **not** a package manager_</sup>
</div>

##

### 🤔 Why?

Setting up the pawn compiler always bothered me. Now, with **pawnctl** it's just a matter of initializing the project and running `pawnctl c`. 😃 

## 🖥️ Installation
 * Head to the **Actions**, pick the latest one and download **pawnctl-win** from it's artifacts. \
 _You can also download it by clicking <a href="https://github.com/christhefrog/pawnctl/actions/runs/9371365016/artifacts/1567803864">pawnctl-win</a>. (it may not be the latest build)_
 * Unpack **pawnctl.exe** to a directory where you have write permissions, prefferably **Desktop** or **Documents**. Make sure **pawnct.exe** is in it's own folder. \
 (_e.g. Documents\pawnctl\pawnctl.exe_)
 * Add the directory to **Path**.
 * Run `pawnctl u` to download the latest compiler.

## 🚀 Usage
 * While in your server directory, run `pawnctl i`.
 * Go through the initialization process. _(if you're running open.mp, you may not need to change anything at all 😊)_
 * Compile the project.
    * Use `pawnctl c` to build a debug version of your code. 
    * Use `pawnctl w` to build a debug version of your code every time a file changes. 
    * Add `release` after the command to build a release version of your code. _(e.g. `pawnctl c release`)_

## 📄 Commands
*  `u(pdate) (version)` \
    Download the specified version of pawn compiler or download the latest one if `version` not specified.
*  `i(nit)` \
    Initialize a new pawnctl project.
*  `c(ompile) (profile)` \
    Compile the project with the specified profile or compile it with the default one if `profile` not specified.
*  `w(atch) (profile)` \
    Compile the project with the specified profile every time a file changes.

## 🛠️ Building
Just copy the repo and run: \
`go build main.go -o pawnctl.exe` \
\
The project also has a github action dedicated to compiling the project. 

<!-- Contact -->
## 🤝🏻 Contact
Discord - @christhefrog
