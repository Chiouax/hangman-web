# hangman-web
# Hangman Web

Un jeu du pendu en ligne avec une interface futuriste, développé en Go.

## Table des matières
1. [Installation](#installation)
2. [Lancement](#lancement)
3. [Architecture du projet](#architecture-du-projet)
4. [Routes](#routes)
5. [Équipe](#équipe)

## Installation

### Prérequis
- Go 1.21 ou supérieur
- Un navigateur web moderne

### Étapes d'installation
1. Clonez le dépôt :
```bash
git clone [votre-url-de-repo]
cd hangman-web
```

2. 
hangman-web/
├── main.go
├── static/
│   ├── css/
│   │   └── style.css
│   └── images/
│       ├── hangman-0.png
│       ├── hangman-1.png
│       ├── hangman-2.png
│       ├── hangman-3.png
│       ├── hangman-4.png
│       ├── hangman-5.png
│       └── hangman-6.png
├── templates/
│   ├── home.html
│   ├── game.html
│   ├── end.html
│   └── scores.html
└── README.md
```

## Lancement

Pour lancer le serveur :
```bash
go run main.go
```

Le serveur démarre sur `http://localhost:8080`

## Architecture du projet

### Organisation des fichiers
- `main.go` : Point d'entrée de l'application et logique du serveur
- `static/` : Fichiers statiques (CSS, images)
- `templates/` : Templates HTML
- `scores.json` : Stockage des scores (créé automatiquement)

### Routes

#### Routes de vue
- `/` : Page d'accueil (GET)
- `/game` : Page de jeu (GET)
- `/end` : Page de fin de partie (GET)
- `/scores` : Tableau des scores (GET)

#### Routes de données
- `/start` : Démarrer une nouvelle partie (POST)
- `/play` : Soumettre une lettre/mot (POST)
- `/static/*` : Servir les fichiers statiques

## Fonctionnalités

1. **Interface utilisateur**
   - Design futuriste avec effets néon
   - Responsive design
   - Animations fluides

2. **Gameplay**
   - Trois niveaux de difficulté
   - Système de points de vie
   - Visualisation du pendu
   - Historique des lettres essayées

3. **Système de scores**
   - Persistance des scores
   - Tableau des meilleurs scores
   - Statistiques par difficulté

## Développement

### Phases clés
1. Mise en place de l'architecture de base
2. Implémentation de la logique de jeu
3. Développement de l'interface utilisateur
4. Tests et corrections de bugs
5. Finalisation et optimisations

### Gestion du temps
1. **Semaine 1**
   - Architecture du projet
   - Routes de base
   - Templates HTML

2. **Semaine 2**
   - Logique de jeu
   - Gestion des sessions
   - Système de scores

3. **Semaine 3**
   - Design et CSS
   - Tests et débogage
   - Optimisations

### Stratégie de documentation
- Documentation Go officielle
- Tutoriels en ligne pour Go web
- Ressources CSS modernes
- Stack Overflow pour les problèmes spécifiques

## Équipe

  - Développement backend
  - Logique de jeu
  - Gestion du projet
