# TODO List : Tests du Module Réseau Serveur (TDD)

## 1. Gestion des Connexions TCP de Base
- [ ] Le serveur accepte une nouvelle connexion TCP
- [ ] Le serveur enregistre la connexion avec son adresse IP et un horodatage précis dans les logs
- [ ] Le serveur gère la déconnexion volontaire d'un client (EOF) sans crasher
- [ ] Le serveur gère la déconnexion abrupte (perte de réseau) sans crasher
- [ ] Le serveur logue la déconnexion du client

## 2. Lecture et Découpage du Flux (Framing)
- [ ] Le serveur lit les messages encodés en UTF-8
- [ ] Le serveur utilise correctement le saut de ligne (`\n` ou `0x0A`) comme délimiteur de message
- [ ] **Fragmentation :** Si un message arrive découpé en plusieurs paquets TCP, le serveur le met en buffer et attend le `\n` avant de le traiter
- [ ] **Coalescence :** Si plusieurs commandes arrivent dans un seul paquet TCP (ex: `cmd1\ncmd2\n`), le serveur les traite séparément l'une après l'autre

## 3. Sécurité et Limites du Réseau
- [ ] Le serveur rejette ou nettoie les caractères de contrôle indésirables envoyés par un client pour éviter les manipulations de terminal
- [ ] Le serveur rejette les messages dépassant la limite de taille (1024 octets par ligne) pour éviter les attaques par saturation de mémoire
- [ ] Un client qui spamme des commandes ou ouvre des connexions trop rapidement est détecté et logué (prévention des abus)

## 4. Écriture et Concurrence (Goroutines & Channels)
- [ ] Le serveur peut envoyer un message à un client spécifique via sa connexion.
- [ ] Si deux clients envoient des données exactement à la même milliseconde (Data Race), le serveur traite les deux sans corrompre sa mémoire.
- [ ] Le serveur peut diffuser (broadcast) un message à plusieurs clients connectés.
- [ ] **Test de robustesse :** La diffusion de messages ne s'interrompt pas et ne fait pas planter le serveur si un des destinataires se déconnecte au moment précis de l'envoi