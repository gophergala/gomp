#!/bin/bash
OLD_PATH=$(pwd)
NEW_PATH=$GOPATH/src/github.com/gophergala/gomp/tests

cd $NEW_PATH


for entry in $(ls *.go)
do
	if [ $entry != 'run_tests.sh' ]; then
		echo -n 'Processing '$entry'... '
		OLD_SOURCE=$entry
		OLD_PROG=${OLD_SOURCE%'.go'}
		OLD_RESULT=$OLD_PROG'_result'
		NEW_SOURCE=${OLD_SOURCE%'.go'}'_modified.go'
		NEW_PROG=${NEW_SOURCE%'.go'}
		NEW_RESULT=$NEW_PROG'_result'
		go build $OLD_SOURCE
		gompp < $OLD_SOURCE > $NEW_SOURCE
		go build $NEW_SOURCE
		./$OLD_PROG | sort > $OLD_RESULT
		./$NEW_PROG | sort > $NEW_RESULT

		diff $OLD_RESULT $NEW_RESULT > tmp

		if [$(cat tmp) == '']; then
			echo 'Passed'
		else
			echo 'Failed'
			exit 1
		fi

		rm -rf $NEW_SOURCE $NEW_PROG $NEW_RESULT $OLD_PROG $OLD_RESULT tmp
	fi
done

cd $OLD_PATH
