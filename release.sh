#!/bin/bash

ORG=justone
NAME=pod-justone-archiver
ARCHS="darwin/amd64 linux/amd64 windows/amd64"

set -ex

if [[ ! $(type -P gox) ]]; then
    echo "Error: gox not found."
    echo "To fix: run 'go get github.com/mitchellh/gox', and/or add \$GOPATH/bin to \$PATH"
    exit 1
fi

if [[ -z $GITHUB_TOKEN ]]; then
    echo "Error: GITHUB_TOKEN not set."
    exit 1
fi

if [[ ! $(type -P github-release) ]]; then
    echo "Error: github-release not found."
    exit 1
fi

VER=$1

if [[ -z $VER ]]; then
    echo "Need to specify version."
    exit 1
fi

PRE_ARG=
if [[ $VER =~ pre ]]; then
    PRE_ARG="--pre-release"
fi

git tag $VER

echo "Building $VER"
echo

rm -v ${NAME}* || true
gox -ldflags "-X main.version=$VER" -osarch="$ARCHS"

# Create zip archive of each binary
for file in ${NAME}_*; do
    arch=${file#"${NAME}_"}
    arch_no_ext=${arch%%.*}
    arch_ext=${arch#"$arch_no_ext"}

    final_bin=$NAME$arch_ext
    mv -v $file $final_bin
    zip $NAME-$VER-$arch_no_ext.zip $final_bin
    rm $final_bin
done

echo "* " > desc
echo "" >> desc

echo "$ sha1sum ${NAME}-*" >> desc
sha1sum ${NAME}-* >> desc
echo "$ sha256sum ${NAME}-*" >> desc
sha256sum ${NAME}-* >> desc

vi desc

git push --tags

sleep 2

cat desc | github-release release $PRE_ARG --user ${ORG} --repo ${NAME} --tag $VER --name $VER --description -
for file in ${NAME}-*; do
    github-release upload --user ${ORG} --repo ${NAME} --tag $VER --name $file --file $file
done
