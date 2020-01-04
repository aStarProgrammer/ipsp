# ipsp
IPSP (Inter Planet Site Publisher) is a tool to publish site to ipfs network

## Background

IPFS (Inter Planet File System [IPFS](https://ipfs.io)) is a peer-to-peer hyperlink protocol which is used to publish content. We can publish a web site  on IPFS as we publish a site on http.

But as IPFS is an p2p system, file published on IPFS cannot be changed, if we changed a file and publish to IPFS again, it is a completely new file from the old one.  Changing files of a IPFS file is not encouraged. So generally sites that are built on ASP.NET Java PHP which have a lot of scripts are not the best option when you want to publish a site to IPFS. Static website based on HTML and CSS is the best option.

IPCS is the tool to create static html site that you can publish to IPFS.

The site created by IPCS looks as follows:

- Site Root Folder
  - index.html
  - more1.html
  - more2.html
  - Pages
    - A1_xxxxxx.html
    - A2_xxxxxx.html

IPSP is the tool to monitor this site periodically and if the site changes, publish the site to ipfs and ipns

## Install

Download  release according to your platform and unzip it.


## Build

If you can not find a release for your platform, build it from source code as follows:

1. Install go

2. Install git

      	Download and install
      		https://git-scm.com/download
      	OR
      		sudo apt-get install git	

3. Install mingw(Windows)

4. Install Liteide (https://github.com/visualfc/liteide)


   ​	*Windows/Linux/MacOSX just download and install

   ​	*Raspbian

   ​		Download source (qt4 Linux 64)code and compile as follows:

  

   ```bash
       sudo apt-get update
       sudo apt-get upgrade
       sudo apt-get install git
       git clone https://github.com/visualfc/liteide.git
       sudo apt-get install qt4-dev-tools libqt4-dev libqtcore4 libqtgui4 libqtwebkit-dev g++
       cd liteide/build
       ./update_pkg.sh
       export QTDIR=/usr
       ./build_linux.sh
       cd ~/liteide/liteidex
       ./linux_deploy.sh
       cd ~/liteide/liteidex/liteide/bin 
      ./liteide
   ```

5. Open ipsp with liteide 

6. Select the platform you needed, modify current environment according to step 1 and 3
   Modify GOROOT and PATH

7. Compile->Build

## Usage

Firstly, you need to build an ipfs service on your server that you want to use ipsp.

You can do it according to following link:

https://docs.ipfs.io/



ipsp has only 1 command:


* Start Monitor
  	
```bash
ipsp -SiteFolder -MonitorInterval
```

​	Monitor the site with SiteFolder every MonitorInterval (second)

​	Example:

```bash
ipsp -SiteFolder "F:\TestSite" -MonitorInterval 600
```

​	 ipsp will monitor site folder F:\TestSite every 600 seconds (10 minutes), if the site changed, publish it to ipfs and ipns again.

​	

```
Note: Actually, the site will try to publish the site to ipfs every 10 minutes , and check the QmID of site folder, if the anything from the site folder changed, the QmID returned from ipfs will change, and ipsp will try to publish this site to ipns again. If nothing changed in the site folder, QmID will not change.
```



## Raise A Issue

Send email to sdxianchao@gmail.com 

## Maintainers

[@aStarProgrammer](https://github.com/aStarProgrammer).


## License

[MIT](LICENSE)