#!/bin/awk -f

BEGIN {
	found = 0;
	ZARGV[1]=ARGV[1];
	ZARGV[2]=ARGV[2];
	ARGV[1]="";
	ARGV[2]="";
}
$0 ~ ZARGV[1] {
	if (!found) {
		found = 1;
		$0 = substr($0, index($0, ZARGV[1]) + length(ZARGV[1]));
	}
}
$0 ~ ZARGV[2] {
	if (found) {
		found = 2;
		$0 = substr($0, 0, index($0, ZARGV[2]) - length(ZARGV[2]));
	}
}
	{ if (found) {
		print;
		if (found == 2)
			found = 0;
	}
}
