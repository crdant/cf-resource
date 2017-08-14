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
    -p | --package-name )
      packageName=$2
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

export GOPATH=$PWD
srcDir="src/$packageName"
mkdir -p $srcDir
cp -R ./$inputDir/* $srcDir/*
cp $srcDir

glide install
go build -o $outputDir/check ./check/cmd/check
go build -o $outputDir/in ./in/cmd/in
go build -o $outputDir/out ./in/cmd/out
