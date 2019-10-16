# PRR - Laboratoire 1 : Synchronisation d'horloge

## Utilisation

Pour utiliser les horloges facilement et à leurs pleines capacités, l'utilisation de docker est conseillée. Ainsi, les horloges fonctionneront un peu plus comme si elles tournaient sur des machines différentes.

### Sur la machine hôte

Si l'utilisation de Docker n'est pas souhaitable dans votre cas, voici comment lancer une horloge maîtresse et une horloge esclave sur la machine hôte.

Premièrement, si vous désirez utiliser un IDE tel que GoLand, il vous suffira d'ouvrir le dossier `src` et de laisser l'IDE télécharger les dépendances. Lancez ensuite la fonction `main()` dans `master.go` et dans `slave.go` (dans l'ordre que vous désirez).

Si vous préférez effectuer ces actions à travers la ligne de commande, les commandes suivantes devraient lancer le maître, puis l'esclave :

```
go build github.com/Laykel/PRR-Lab1/master
go build github.com/Laykel/PRR-Lab1/slave

./master
./slave
```

### Avec Docker

Afin de lancer le maître, il suffit d'entrer la commande suivante, après vous être placé dans le répertoire `src/github.com/Laykel/PRR-Lab1`.

```
docker-compose up -d clock-master
```

Cela téléchargera le container `golang` avant de copier les fichiers sources dans le container, le construire et le lancer en arrière-plan.

Ensuite, pour lancer les esclaves et afficher les logs, entrer cette commande :

```
docker-compose up --scale clock-slave=5
```

Si vous préférez lancer les containers en arrière plan, ajouter le flag `-d` et utilisez ensuite la commande

```
docker-compose logs
```

pour voir la sortie des programmes.

## Fonctionnement

## Problèmes connus

- Dans le slave, pas de timeout lors de l'attente de DELAY_RESPONSE.

- Nous n'avons pas trouvé la source d'une erreur arrivant de manière peu reproductible "unexpected EOF".
