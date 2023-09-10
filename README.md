# 🔎 naldump

Tool for parsing an H.264 or H.265 Network Abstraction Layer (NAL) Annex B
bytestream and printing the size and type of each NAL unit.

Input may be read from a file or standard input. If read from standard input,
the time (in seconds since start) is printed.

H.264 is assumed unless `-hevc` flag is specified. To exclude certain NAL
types, use `-exclude` flag. For instance, to exclude H.264 Supplemental
Enhancement Information (SEI) units, specify `-exclude 6`. To exclude more
than one type, use a comma-separated list.

## Example
```
$ cat myfifo | ./naldump -hevc
0.214892	0	1	12514	TRAIL_R
0.223971	1	1	12373	TRAIL_R
0.238800	2	1	13168	TRAIL_R
0.256814	3	1	12807	TRAIL_R
0.265341	4	1	12755	TRAIL_R
0.278517	5	1	12351	TRAIL_R
0.292916	6	1	13066	TRAIL_R
0.302399	7	1	12795	TRAIL_R
0.309455	8	1	12488	TRAIL_R
0.316250	9	1	12445	TRAIL_R
0.326067	10	1	12690	TRAIL_R
0.339279	11	1	12950	TRAIL_R
0.351582	12	1	13124	TRAIL_R
0.373358	13	1	13034	TRAIL_R
0.447804	14	1	12942	TRAIL_R
0.461442	15	1	12845	TRAIL_R
0.522417	16	1	12644	TRAIL_R
0.570705	17	1	12802	TRAIL_R
0.650681	18	1	12884	TRAIL_R
0.655817	19	1	13071	TRAIL_R
0.657328	20	32	23	video parameter set
0.658980	21	33	34	sequence parameter set
0.661174	22	34	7	picture parameter set
0.741680	24	19	53082	IDR_W_RADL
0.772504	25	1	12088	TRAIL_R
0.848755	26	1	12808	TRAIL_R
0.861298	27	1	13183	TRAIL_R
```

## Build

```
go build
```
