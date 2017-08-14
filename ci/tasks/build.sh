#!/bin/sh

inputDir=  outputDir=

while [ $# -gt 0 ]; do
  case $1 in
    -i | --input-dir )
      inputDir=$2
      shift
      ;;
    -o | --output-dir )
      outputDir=$2
      shift
      ;;
    * )
      echo "Unrecognized option: $1" 1>&2
      exit 1
      ;;
  esac
  shift
done

error_and_exit() {
  echo $1 >&2
  exit 1
}

if [ ! -d "$inputDir" ]; then
  error_and_exit "missing input directory: $inputDir"
fi
if [ ! -d "$outputDir" ]; then
  error_and_exit "missing output directory: $outputDir"
fi
if [ -z "$packageName" ]; then
  error_and_exit "missing package name (from task parameters)"
fi

mkdir gopath
export GOPATH=$PWD/gopath
goPathSrc="${GOPATH}/src/$packageName"
mkdir -p $goPathSrc
cp -r $inputDir/* $goPathSrc
cd $goPathSrc

glide install
go build -o $outputDir/check ./check/cmd/check
go build -o $outputDir/in ./in/cmd/in
go build -o $outputDir/out ./out/cmd/out
