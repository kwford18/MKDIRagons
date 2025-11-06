# MKDIRagons
A tool to help in D&D character creation.  
Build a basic character using the blank TOML file, which is then parsed into a feature-complete 5e character!  
  
Data is fetched from the 5e API, and in future updates missing data will be scraped off 5e Wikidot  
  
# Currently supported commands:  
- build: Build a JSON character file using the provided TOML file  
- load:  Load the provided JSON character and print their information  
- empty: Generate an empty TOML template file for making characters in toml-empty/  
  
The build & load commands both support the --file/-f flags, which are  
for providing a path to the file to build from/load respectively. Build looks for the **toml-characters** directory by default, but this can be overriden with the -f flag.  
Run MKDIRagons empty to generate a template file in this directory.  

The build command also supports the following two switches:  
- --print/-p:  Print the content of the built file (equivalent to loading
  file)  
- --rollHP/-r: Use random number generation for calculating health
  instead of the hit die average  