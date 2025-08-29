# WordMon-Go 🎮

Un jeu de capture de mots inspiré de Pokémon, développé en Go.

[![CI](https://github.com/jusgaga/wordmon-go/workflows/CI/badge.svg)](https://github.com/jusgaga/wordmon-go/actions)

## 🎯 Description

WordMon-Go est un jeu où les joueurs capturent des "créatures-mots" en résolvant des défis linguistiques. Chaque mot a une rareté et des points d'expérience, et les joueurs progressent en niveau en accumulant de l'XP.

## 🚀 Fonctionnalités

- **Système de capture** : Capturez des mots en résolvant des anagrammes
- **Progression des joueurs** : Système d'XP et de niveaux
- **Mots rares** : Différents niveaux de rareté (Common, Rare, Legendary)
- **Machine à états** : Gestion des rencontres avec transitions d'état
- **API REST** : Interface HTTP pour interagir avec le jeu
- **Configuration flexible** : Support YAML, TOML et JSON

## 🛠️ Technologies

- **Go 1.25+** : Langage principal
- **Gin** : Framework web
- **PostgreSQL** : Base de données
- **Staticcheck** : Analyse statique du code
- **Delve** : Débogueur Go

## 📋 Prérequis

- Go 1.25 ou supérieur
- PostgreSQL (optionnel pour le développement)
- Make (pour utiliser le Makefile)

## 🚀 Installation

1. **Cloner le repository**

   ```bash
   git clone https://github.com/jusgaga/wordmon-go.git
   cd wordmon-go
   ```

2. **Installer les dépendances**

   ```bash
   go mod download
   ```

3. **Installer les outils de développement**
   ```bash
   make tools
   ```

## 🎮 Utilisation

### Démarrage rapide

```bash
# Compiler le projet
make build

# Lancer le serveur
./bin/server
```

### Commandes Makefile

```bash
# Vérification complète de la qualité
make all

# Formatage du code
make fmt

# Vérification statique
make vet

# Linting
make lint

# Tests
make test

# Tests avec détection de race conditions
make test-race

# Tests avec couverture
make test-cover

# Compilation
make build

# Nettoyage
make clean
```

## 🧪 Tests

### Exécution des tests

```bash
# Tous les tests
go test ./...

# Tests avec détection de race conditions
go test -race ./...

# Tests avec couverture
go test -cover ./...
```

### Couverture minimale

Le projet maintient une couverture de tests d'au moins **60%**.

## 🔍 Qualité du code

### Outils utilisés

- **`go fmt`** : Formatage automatique du code
- **`go vet`** : Vérification statique
- **`staticcheck`** : Analyse statique avancée

### Critères de qualité

- ✅ Code correctement formaté
- ✅ Aucune erreur de vérification statique
- ✅ Respect des règles de linting
- ✅ Tests passants
- ✅ Couverture suffisante

## 🐛 Debugging

### Installation de Delve

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
```

### Utilisation

```bash
# Lancer le serveur en mode debug
dlv debug ./cmd/server

# Breakpoints recommandés
# - Résolution des combats
# - Gestion des captures
# - Transitions d'état des rencontres
```

## 📚 Documentation

### Génération de la documentation

```bash
# Lancer le serveur de documentation
godoc -http=:6060

# Ouvrir http://localhost:6060 dans votre navigateur
```

### Packages documentés

- `internal/core` : Types et logique métier
- `internal/config` : Configuration du jeu
- `internal/api` : Interface HTTP

## 🔄 CI/CD

### GitHub Actions

Le projet utilise GitHub Actions pour la CI avec les étapes suivantes :

1. **Formatage** : Vérification du formatage avec `go fmt`
2. **Vérification statique** : Exécution de `go vet`
3. **Linting** : Analyse avec `staticcheck`
4. **Tests** : Exécution de tous les tests
5. **Race conditions** : Détection des conditions de course
6. **Couverture** : Vérification de la couverture ≥ 60%
7. **Compilation** : Build du binaire

### Déclencheurs

- Push sur `main` et `master`
- Pull requests sur `main` et `master`

## 📁 Structure du projet

```
wordmon-go/
├── cmd/server/          # Point d'entrée principal
├── internal/            # Code interne
│   ├── api/            # Handlers HTTP et serveur
│   ├── config/         # Configuration du jeu
│   └── core/           # Logique métier et types
├── configs/            # Fichiers de configuration
├── db/                 # Migrations de base de données
├── .github/workflows/  # Workflows CI/CD
└── Makefile           # Commandes de build et qualité
```

## 🤝 Contribution

1. Fork le projet
2. Créer une branche feature (`git checkout -b feature/AmazingFeature`)
3. Commit les changements (`git commit -m 'Add some AmazingFeature'`)
4. Push vers la branche (`git push origin feature/AmazingFeature`)
5. Ouvrir une Pull Request

### Standards de code

- Respecter le formatage Go (`go fmt`)
- Ajouter des tests pour les nouvelles fonctionnalités
- Maintenir la couverture de tests ≥ 60%
- Documenter les nouvelles fonctions et types

## 📄 Licence

Ce projet est sous licence MIT. Voir le fichier `LICENSE` pour plus de détails.

## 🆘 Support

Pour toute question ou problème :

1. Vérifier les [Issues](https://github.com/jusgaga/wordmon-go/issues)
2. Créer une nouvelle issue si nécessaire
3. Consulter la documentation avec `godoc`

---

**Développé avec ❤️ en Go**
