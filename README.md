# mpdf

Merge PDF files in a given directory in specified order. No more Adobe Acrobat subscription or uploading sensitive files to random sites.

## Usage
```shell
mpdf [-c] [-r] [-f] [-in dir]  [-out merged]
```

Merge all PDF files and output the merged file in the current directory:
```shell
mpdf
```

By default, `mpdf` merges PDF files it finds in ascending lexicographical order of filename.

`mpdf` will pick a random filename for the merged file, if it is not specified via `-out` flag, or the file path specified by `-out` flag is a directory.

Merge all PDF files in directory `/path/to/pdf/dir`, then output the merged PDF file to specified destination:
```shell
mpdf -in /path/to/pdf/dir -out /path/to/merged.pdf
```

Merge pdf files in the order of last modification time:
```shell
mpdf -c -in /path/to/pdf/dir -out /path/to/merged.pdf

``` 

Use `-r` flag to reverse the merge ordering:
```shell
mpdf -r -in /path/to/pdf/dir -out /path/to/merged.pdf
``` 

By default `mpdf` fails if there already existed a file with specified destination filename. To overwrite, use `-f` flag.

`mpdf` fails if the parent directory of specified destination doesn't exist.

## Third-party Dependencies

* [`pdfcpu`](https://pdfcpu.io) for merging PDF files
