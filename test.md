Tu es un expert Go et CLI design.
Je veux que tu m’écrives un projet Go complet appelé wrapper qui permet de gérer des profils de configuration pour n’importe quel binaire (ex: vault, aws, kubectl).

## Objectif

Pouvoir lancer un binaire via mon **wrapper** en gérant des profils d’environnement.

## Exemple d’usage
```
vault profile list
vault profile create prod
vault profile set prod
vault profile get
vault profile default prod
vault status
```
Ici, vault est en réalité mon wrapper :

```
vault() { wrapper vault "$@"; }
```

## Fonctionnalités attendues
### Détection du mode

Si la commande est profile … → c’est une commande interne du wrapper.

Sinon → exécuter le vrai binaire (vault, aws, etc.) avec le profil actif chargé.

### Gestion des profils

Les profils sont stockés dans ~/.config/wrapper/<bin>/.

Format de stockage : simple fichier .env (clé=valeur).

### Commandes supportée

profile list → liste les profils existants.

profile create <name> → crée un fichier <name>.env.

profile set <name> → définit le profil courant (symlink current.env).

profile get → affiche le profil courant.

profile default <name> → définit le profil par défaut (au cas où aucun profil n’est encore set).

### Exécution du vrai binaire

Charger les variables d’env depuis current.env (ou le profil par défaut).

Puis exécuter le vrai binaire via exec.Command.

Rediriger stdin, stdout, stderr directement (comme un wrapper transparent).

### Résolution du vrai binaire

Le wrapper doit trouver le “vrai” binaire sans se boucler lui-même.

Stratégie : rechercher dans $PATH en ignorant son propre chemin.

Exemple : si le wrapper est lancé en tant que vault, il doit trouver /usr/bin/vault et exécuter celui-ci.

### Ergonomie

Si aucun profil n’est encore défini → message d’erreur clair.

Si un profil par défaut existe → l’utiliser automatiquement.

Les erreurs doivent être joliment formatées.

Possibilité d’utiliser vault() fonction shell comme :

vault()   { wrapper vault "$@"; }
aws()     { wrapper aws "$@"; }
kubectl() { wrapper kubectl "$@"; }

## Exemple d’utilisation

```
# création de profil
vault profile create prod
# → édite ~/.config/wrapper/vault/prod.env

# définir un profil courant
vault profile set prod

# afficher le profil courant
vault profile get
# → "prod"

# exécuter le vrai vault avec ce profil
vault toto
# → lance /usr/bin/vault toto avec les env chargés depuis prod.env
```

### Contraintes techniques

Écrit en Go 1.25+
Faire du DDD

### Mission

Écris-moi le code complet de ce projet wrapper, prêt à compiler (go build -o wrapper .), avec :

Gestion des profils (list, create, set, get, default)

Exécution transparente des binaires avec le profil actif

Recherche correcte du “vrai binaire” dans $PATH