# mpdf

Merge PDF files in a given directory.

## Usage
```shell
mpdf [-crf] [-in dir]  [-out merged]
```

Merge all PDF files and output the merged file in the current directory:
```shell
mpdf
```

By default, `mpdf` merges PDF files it found in ascending lexicographical order of filename.

`mpdf` will pick a random filename for the merged file if it is not specified via `-out` flag. It guarantees the filename never clash with those of existing files in the destination directory.

Merge all pdf files in directory `/path/to/pdf/dir`, then output the merged PDF file to specified destination:
```shell
mpdf -in /path/to/pdf/dir -out /path/to/merged/pdf
```

Merge pdf files in the order of last modification time:
```shell
mpdf -c -in /path/to/pdf/dir -out /path/to/merged/pdf

``` 

Use `-r` flag to reverse the merge ordering:
```shell
mpdf -r -in /path/to/pdf/dir -out /path/to/merged/pdf
``` 

By default `mpdf` fails if there already existed a file with specified destination filename. To overwrite, use `-f` flag.

## Third-party Dependencies

* `pdfcpu` for merging PDF files
