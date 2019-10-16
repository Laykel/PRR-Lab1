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

### Maître

Le maître envoie des messages de type SYNC et FOLLOW_UP à intervalles réguliers (toutes les 2 secondes), cela en mode multicast (avec UDP). Le message SYNC envoie un ID, alors que le FOLLOW_UP envoie un ID et un temps tMaster. En plus de cette tâche, il attend des messages de type DELAY_REQUEST que sont censés lui envoyer les esclaves. Lorsqu'il reçoit un message DELAY_REQUEST, il répond à l'esclave en question avec un message de type DELAY_RESPONSE contenant un ID, et un temps tM.

### Esclave(s)

Les esclaves écoutent les messages multicast envoyé par le maître. Une fois qu'ils ont reçu un SYNC, ils suivent ensuite une séquence SYNC-FOLLOW_UP-DELAY_REQUEST-DELAY_RESPONSE. Avec les paramètres de ces messages, ils déterminent leur propre décalage.

### Divers

Les temps envoyés et traités sont des timestamps en microsecondes.

## Problèmes connus

- Dans le slave, pas de timeout lors de l'attente de DELAY_RESPONSE.

- Nous n'avons pas trouvé la source d'une erreur arrivant de manière peu reproductible "unexpected EOF".
