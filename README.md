# Formule 1 IA04

IA04 - Groupe 1C
Yannis Brena--Label, Adam Hafiz, Hugo Milair, Damien Vaurs

Semestre A23, supervisé par S. Lagrue et H. Willot.

## Préambule

Ce fichier s'inscrit dans le cadre du travail réalisé par le groupe 1C pour le projet de modélisation et de simulation d'un système multi-agent. Il décrit à la fois le backend (Go) et le frontend (React) que le lecteur pourra trouver dans les répertoires GitLab suivant :

- frontend - https://gitlab.utc.fr/ybrenala/formule-1-ia04-front 
- backend - https://gitlab.utc.fr/vaursdam/formule-1-ia04 

## Lancement du projet

### Installation

Les deux parties du projet sont téléchargeables avec les commandes suivantes :

    git clone https://gitlab.utc.fr/ybrenala/formule-1-ia04-front
    git clone https://gitlab.utc.fr/vaursdam/formule-1-ia04 

**Après le clone**, pour le frontend, il se peut que des dépendances doivent être installées. Lors de nos essais, il s'est avéré que la commande suivante était l'unique à réaliser pour pouvoir lancer correctement le projet : 

    npm i react-chartjs-2

**Si l'utilisateur ne clone pas le projet mais a obtenu le projet via le dépôt IA04**, il faut installer toutes les dépendances:

- se placer dans app-react
- réaliser la commande *npm install* pour installer les dépendances.

### Lancement des programmes

En ligne de commande, l'interface utilisateur se lance depuis le sous-repertoire *app-react* avec **npm**:

    npm run start

Quant au projet Go, l'utilisateur peut soit lancer le programme avec :

    go run cmd/launch-simulation.go

Ou en installant dans un premier temps le fichier exécutable :

    go install cmd/launch-simulation.go

L'utilisateur pourra alors exécuter le fichier depuis son répertoire Go.

### Remarques et conseils

Pour le bon fonctionnement de l'interface utilisateur, il est impératif que le backend soit lancé. Aussi, l'utilisateur doit s'assurer que son **port 8080** soit libre pour que les requêtes du frontend atteignent bien le backend.

Aussi, lors de l'utilisation de l'interface utilisateur, il est possible que les graphes s'affichent mal dans l'onglet principal. On conseille à l'utilisateur de cliquer sur le bouton *Simuler un seul championnat*, aller sur un autre onglet, puis revenir sur l'onglet principal. Les graphes s'afficheront alors correctement lors de la simulation.

Enfin, le dossier **python_plots** contient le scripts Python permettant de tracer des graphes intéressants dans le cadre de la simulation mais que nous n'avons pas eu le temps d'inétgrer à l'interface utilisateur. L'utilisateur peut visualiser les dits graphes dans le même dossier. Des indications précises sur l'utilisation de ce script sont disponibles en tête du fichier Python.

## Description du projet

### Objectif du projet

Ce projet a pour but de répondre à la problématique **Quel est le meilleur profil d’un pilote pour obtenir le plus de points ?**. Y répondre permettrait notamment, dans la position d'une équipe, d'éclairer la stratégie de recrutement des pilotes à la lumière de leur personnalité.

Une personnalité est définie par 4 traits de caractère qui influencent le comportement du pilote en course. Ces traits sont les suivants :

- **Aggressivité** : détermine la propension du pilote à tenter des dépassements
- **Concentration** : détermine la capacité du pilote à se concentrer sur la course
- **Confiance** : détermine le niveau de confiance en soi du pilote
- **Docilité** : détermine la docilité du pilote face aux consignes de l'équipe

Les deux premiers de ces traits, une fois fixés, ne varient pas au cours de la carrière du pilote. Les deux derniers, en revanche, peuvent évoluer en fonction des performances en courses des pilotes.

L'objectif de la simulation est ainsi de **déceler quelle personnalité est la plus performante en course**. Pour cela, l'utilisateur peut, à travers l'interface, simuler des championnats ou des courses de Formule 1 et visualiser les résultats des pilotes, des équipes et des différents profils de personnalité. Il peut également modifier les personnalités des pilotes au cours de la simulation pour observer l'impact de ces changements sur les résultats du pilote.
On lui conseille par ailleurs d'observer les pilotes surperformant dans les championnats (ex. le pilote Sargeant, bien que conduisant une mauvaise voiture et ayant un niveau intrinsèque faible, est régulièrement dans les meilleurs pilotes de nos simulations) et d'appliquer leurs personnalités à des pilotes mal classés. Le pilote aura tendance à gagner plus de points et l'utilisateur pourra observer l'impact de la personnalité sur les résultats.

Remarquons que la présentation du projet face aux professeurs ainsi que le graphe contenu dans **python_plots** ont mené à la conclusion que le pilote idéal est **peu agressif**, **peu concentré** et **peu docile**. Il est par ailleurs **très confiant**.

### Modélisation

Ce projet prend en comptes plusieurs éléments des championnats de Formule 1. Les plus importants sont:

- les **pilotes**, les agents de la simulation, qui sont caractérisés par leur personnalité, leur niveau et leur équipe. Les traits de personnalité initiaux de chaccun des 20 pilotes ont été estimés par un suiveur assidu de la Formule 1.
- les **circuits** sur lesquels évoluent les pilotes. 12 circuits régulièrement impliqués dans les championnats ont été modélisés.
- les **courses**, qui constituent le point d'intérêt des simulations.

Au sein des courses, certains éléments sont pris en compte :

- la **météo**, dont la distribution varie en fonction de la situation géographique du circuit.
- les **pneus**, qui s'ils ne sont pas remplacés risque de mener à une crevaison. L'état des pneus influe aussi sur la vitesse des pilotes
- les **arrêts au stand** qui permettent de changer de pneus.

D'autres éléments, notamment les qualifications et essais libres, pourraient être ajoutés. Aussi le modèle pourrait être amélioré pour correspondre au mieux à la réalité.

## Captures d'écran de l'interface graphique

### Simulation de championnats

#### Statistiques pilotes

![Alt text](doc/screens/drivers.png)

#### Statistiques équipes

![Alt text](doc/screens/teams.png)

#### Statistiques personnalités

![Alt text](doc/screens/personnality.png)

### Simulation d'une course

#### Résultats d'une course

![Alt text](doc/screens/race.png)

#### Personnalités d'une course

![Alt text](doc/screens/persoRace.png)

### Modification des personnalités

![Alt text](doc/screens/perso.png)