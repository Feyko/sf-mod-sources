# Satisfactory Mod Sources

Little utility to get all mod sources from https://ficsit.app  
Takes a filepath as an argument. Output is JSON  
No filtering at all on mods, truly gets all the sources, outdated or not!  

If you want to build this for whatever reason, run `go generate -x gql/gen.go` first :slightly_smiling_face:
Also, the generation script uses `sh`, which will not work on Windows system
