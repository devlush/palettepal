# palettepal

### Motivation
A color system for the NES platform is devised using unique a interlacing technique on NTSC signals.  However, the color combinations that yield decent picture quality are rare.  This golang program uses the monte carlo method to seek out viable color combinations in a vast search space (approximately 10^80)


### Architecture
Any number of hosts each run a copy of the program, which in turn spins up golang threads (channels) that search for phase pairs (color combos) that match a predefined filter criteria.  When an interesting result is found, the program will connect to a postgres database and store it.
Deployment of these worker programs is accomplished using Ansible from a controller host.


### Usage
Docker compose files are included for running the postgres database on the controller host.

```
[devlush@beowulf palettepal]$ cd ./docker
[devlush@beowulf palettepal]$ docker-compose

```



This program is deployed using ansible.  Variables such as worker thread count and iterations can be configured in the `inventory.yml` file.  Use the launcher script to kick off the ansible playbook:

```
[devlush@beowulf palettepal]$ ./launch 

PLAY [all] ***************************************

TASK [Gathering Facts] ***************************

```
