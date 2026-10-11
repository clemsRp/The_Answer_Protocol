Audit technique · RFC 42TAP
The Answer Protocol
10 octobre 2026
branche main · commit 3f31f21
Go · 19 266 lignes · 290 commits
auteurs : crappo, aluslu
En l’état, le projet ne passerait pas l’évaluation.
Une section obligatoire du README est vide, et un seul paquet TCP suffit à arrêter le serveur. Le build passe, le monde respecte les minimums et les commandes de base suivent la RFC.

86
problèmes relevés
4
critiques
35
majeurs
47
mineurs
Les quatre points critiques
B1
Critique
Un paquet arrête le serveur
QUIT supprime le joueur mais garde sa session. La commande suivante du même paquet déréférence un joueur nil et le moteur panique.

C: CONNECT alice
C: QUIT
C: LOOK
S: OK connected
S: OK bye
Fatal crash : Engine panic: nil pointer dereference
B2
Critique
Le serveur se fige pour tout le monde
Le moteur et le serveur s’écrivent l’un à l’autre par des canaux bornés. Quand les deux sont pleins, chacun attend l’autre sans fin. 165 connexions en rafale suffisent.

C: WHO
   (aucune réponse)
nouvelle connexion
S: ERR 901 SERVER_FULL
moteur bloqué  engine/utils.go:77
serveur bloqué server/server.go:61
B3
Critique
Ctrl+C ne l’arrête plus
Une fois figé, le serveur affiche « Server and Engine stopped with success » mais le processus reste vivant : s.Stop() attend des goroutines bloquées. Il faut un kill -9.

R1
Critique
README : section obligatoire vide
« Group Contributions » contient deux lignes # TODO. La grille arrête l’évaluation dès la rubrique README.

Répartition par catégorie
Bugs de fonctionnement
20
Sécurité et robustesse
11
Protocole, journalisation, clients
21
Architecture et qualité du code
12
Build, Makefile et dépôt
9
Tests et README
13
Critique
Majeur
Mineur
Verdict par rubrique de la grille
Les 17 rubriques de la feuille d’évaluation, dans l’ordre.

#	Rubrique	Verdict	Motif principal
1	Content (livrables)	Oui	Makefile, serveur, CLI, GUI, world.json, README présents ; Go uniquement
2	README Documentation	Non	R1 : la grille arrête l’évaluation ici
3	World requirements	Oui	12 salles en boucles, 3 rôles, 11 objets, 4 quêtes de 2 types
4	Build	Oui	Toutes les cibles passent ; attention à D1
5	Protocol Compliance	Non	P1 à P5 ; B1 sur une entrée inhabituelle
6	Server Behaviour	Partiel	Diffusion ciblée correcte ; validation du monde incomplète (B20)
7	CLI Client	Non vérifié	Non lancé en interactif
8	GUI Client	À risque	G1 : pas de vue journal (lecture du code)
9	Robustness	Non	B2 sur connexions rapides ; B4 à B6 sur les groupes
10	Network Features	Non	S1 et S2 : caractères de contrôle et UTF-8 invalide acceptés
11	Inventaire et PNJ	Partiel	Commandes fonctionnelles ; formats P1 à P3, P7
12	Data Integrity	Partiel	JSON valide, état cohérent ; identifiants de salle incohérents (P10)
13	Dynamic Item Management	Oui	Le parcours de la grille passe ; B9 hors parcours
14	Combat System	Partiel	PV, riposte, réapparition et ERR 405 corrects ; R3, B12
15	Quest System	Partiel	2 types et ERR 406 corrects ; B8, R2
16	Server Logging	Partiel	JSON, niveaux et horodatage corrects ; L1, L2
17	Recode Exercise	À risque	D1 : une modification de world.json n’est pas recompilée
Le drapeau « Crash » de la grille s’applique tant que B1 n’est pas corrigé.

Tous les constats
Bugs de fonctionnement (20)
Sécurité et robustesse (11)
Conformité au protocole (13)
Journalisation (4)
Monde et clients (4)
Architecture et qualité du code (12)
Build, Makefile et dépôt (9)
Tests (4)
README (9)
Tous
Critiques
Majeurs
Mineurs
Reproduits
Lecture du code
Chercher : QUIT, groupe, README…
86 / 86
Bugs de fonctionnement
20 constats
Trois défauts arrêtent ou figent le serveur depuis un simple client TCP. Les autres faussent les groupes, les quêtes et les combats.

Arrêt et blocage
B1
QUIT suivi d’une commande arrête le serveur
Critique
QUIT supprime le joueur mais garde la session ; la commande suivante du même paquet déréférence un joueur nil.
ReproduitCONNECT alice\nQUIT\nLOOK\n → Fatal crash : Engine panic: nil pointer dereference, processus terminé, tous les joueurs déconnectés.
engine/command.go:41-42 · engine/handle_commands.go:46
B2
Interblocage serveur–moteur
Critique
Le moteur écrit dans ServerOutput pendant que le serveur écrit dans JoinChan ou ServerInput. Quand les deux canaux sont pleins, chacun attend l’autre.
Reproduit165 cycles « connexion, CONNECT, fermeture » en moins d’une seconde ; ou 4 clients envoyant 10 000 NOTIFY_PLAYER_POSITION chacun. WHO ne répond plus, toute connexion reçoit ERR 901 SERVER_FULL. En jeu normal (6 joueurs à 60 Hz pendant 12 s), le blocage ne s’est pas produit.
server/server.go:52-66 · engine/utils.go:22-83
B3
Arrêt impossible après blocage
Critique
Ctrl+C affiche « stopped with success » mais s.Stop() attend des goroutines bloquées.
ReproduitProcessus toujours vivant 3 s après SIGINT ; kill -9 nécessaire.
server/server.go:158-169 · cmd/server/main.go:95
Groupes
B4
Prise de contrôle d’un groupe
Majeur
Un joueur promu puis exclu garde son drapeau promotion. Dans un autre groupe, GROUP ACCEPT le rend chef.
Reproduitbob, jamais promu dans le groupe de carol, obtient OK new_leader=bob puis exclut carol.
engine/group.go:90-125, 259-283
B5
Promotion verrouillée
Majeur
Si le joueur promu est exclu ou part, send_promotion du chef n’est jamais remis à zéro.
ReproduitGROUP PROMOTE dave → ERR 411 NO_PERMISSION pour le chef légitime.
engine/group.go:236, 252
B6
Un chef peut s’exclure lui-même
Majeur
Le groupe reste sans chef : personne ne peut inviter ni promouvoir.
ReproduitGROUP KICK eve par eve → OK ; fred reçoit ensuite ERR 411 sur INVITE et PROMOTE.
engine/group.go:90-125
B7
Diffusion aux invités erronée
Mineur
inform_group_invitations teste les invitations de l’émetteur au lieu de celles du destinataire.
LectureL’événement part vers tout le monde ou vers personne.
engine/utils.go:63-72
Quêtes et objets
B8
La progression d’une quête ne bouge jamais
Majeur
refreshQuestProgress ne traite que le statut available, alors qu’une quête acceptée passe en progress.
ReproduitMaelis vaincue, QUESTS affiche encore "progress":"0/1". Seul COMPLETE_QUEST la passe à 1/1.
engine/command.go:440, 563-566
B9
Duplication d’objets
Majeur
Les récompenses ajoutent la définition partagée de l’objet à chaque joueur. Le sujet exige des instances uniques.
ReproduitDeux joueurs battent Bob, chacun reçoit cancoillotte. Après DROP, LOOK renvoie "items":["cancoillotte","cancoillotte"].
engine/combat_system.go:403-407 · engine/command.go:521-525
B10
USE fait disparaître les armes
Mineur
Équiper une arme la retire de l’inventaire et détruit l’arme précédente.
ReproduitUSE pioche → INVENTORY renvoie [].
engine/command.go:741-765
B11
TAKE ignore le champ obtainable
Mineur
Sans effet aujourd’hui : les 11 objets valent true.
LectureAucun test de Obtainable dans handleCmdTake.
engine/command.go:349-375
Combat
B12
MOVE reste autorisé en combat
Majeur
Le joueur attaque depuis une autre salle, le combat s’annule sans FLEE et le PNJ retrouve tous ses points de vie.
ReproduitMOVE east puis ATTACK maelis → "combat_state":"CANCELLED".
engine/command.go:207-224
B13
Déconnexion en combat : session fantôme
Majeur
Le PNJ continue de frapper le joueur absent à chaque délai de 100 s. En groupe, les alliés attendent 100 s à chaque tour du fantôme.
Reproduit100 s après la déconnexion, le log montre EVT COMBAT ALLY_TURN ghost et une attaque de Maelis.
engine/command.go:41-48 · engine/handle_commands.go:123-156
B14
Première frappe du PNJ absente de la réponse
Mineur
Quand le PNJ est plus rapide, sa première attaque n’est pas listée.
ReproduitPremier ATTACK maelis : les PV passent de 100 à 88 sans action correspondante.
engine/combat_system.go:183-201
B15
PV de PNJ négatifs
Mineur
Les dégâts ne sont pas bornés à zéro.
Reproduit"target_hp":-6.
engine/npcs.go:28-31
B16
Butin réservé au chef
Mineur
S’il est absent du combat, personne ne reçoit le butin.
LectureCondition is_leader dans end_combat.
engine/combat_system.go:400-409
B17
Joueur ajouté à un second combat
Mineur
Un joueur déjà en combat est enrôlé quand un allié attaque dans la même salle ; la première session reste orpheline.
LectureaddPlayerToCombat écrase CombatId sans vérification.
engine/combat_system.go:145-153
Chargement du monde
B18
Boucle infinie sur un JSON invalide
Majeur
Une erreur de syntaxe dans world.json empêche le serveur de démarrer, sans message.
Reproduit[{"rooms" {}}] → aucun retour après 3 s.
engine/parser.go:114-123
B19
Clés dupliquées jamais signalées
Mineur
L’erreur de checkDuplicates est ignorée.
ReproduitDeux clés rooms → err=<nil>.
engine/parser.go:116-122
B20
Références non validées
Mineur
Récompense et cibles de quête, butin de PNJ, présence de la salle de départ.
ReproduitMonde sans place_du_village et quête vers des objets inexistants acceptés.
engine/data_models.go:3-11 · engine/npcs.go:17 · engine/player.go:102
Sécurité et robustesse
11 constats
Ni les noms ni les messages ne sont filtrés, et le limiteur de débit laisse passer quatre commandes sans limite.

S1
Caractères de contrôle relayés tels quels
Majeur
ESC, NUL, BEL et CR passent dans les noms et le chat. Un joueur peut effacer ou réécrire le terminal des autres. La RFC §9.2 demande de les rejeter ou de les neutraliser.
ReproduitCONNECT \x1b[31mroot\x1b[0m\x07 accepté, puis diffusé dans EVT ROOM PRESENCE ENTER.
engine/command.go:13, 226
S2
UTF-8 invalide accepté
Majeur
La RFC impose l’UTF-8 pour tous les messages.
ReproduitCONNECT \xff\xfe\xfa → OK connected, octets rediffusés bruts.
engine/command.go:22
S3
Aucune validation du nom
Majeur
Longueur et alphabet libres.
ReproduitNom de 1 000 caractères accepté, puis répété dans chaque LOOK et chaque événement.
engine/command.go:13-39
S4
Limiteur de débit contourné
Majeur
NOTIFY_PLAYER_POSITION, GET_PLAYER_POSITIONS, NOTIFY_ITEM_POSITION et GET_ITEM_POSITIONS ne consomment aucun jeton. C’est la porte d’entrée de B2.
Reproduit40 000 lignes envoyées par 4 clients sans aucun ERR 902.
server/rate_limiter.go:31-38
S5
Accès concurrent à la map des clients
Majeur
s.clients est lue par la goroutine lecteur et écrite par broadcaster. En Go, cela peut arrêter le processus.
Reproduitgo test -race échoue sur TestAntiFloodProtection (DATA RACE).
server/client.go:91 · server/server.go:79
S6
Positions d’objets sans contrôle
Majeur
NOTIFY_ITEM_POSITION accepte n’importe quel nom. La map grossit sans borne et GET_ITEM_POSITIONS renvoie les objets de toutes les salles.
LectureAucune vérification d’existence ni de salle.
engine/command.go:842-875, 928-938
S7
Connexions en rafale ni limitées ni signalées
Mineur
Pas de plafond par IP, pas de log WARN. Le sujet demande de surveiller les « rapid connections ».
Reproduit2 000 connexions depuis une IP : seules des lignes INFO.
server/server.go:39-48
S8
Ligne trop longue : coupure sans réponse
Mineur
Au-delà de 1 024 octets, la connexion est fermée sans ERR.
ReproduitLigne de 2 000 octets → reset TCP.
server/client.go:66-74
S9
Une panique du moteur arrête tout
Mineur
recover renvoie l’erreur à main, qui stoppe serveur et moteur. Aucune commande n’est isolée.
ReproduitDémontré par B1.
engine/engine.go:52-57 · cmd/server/main.go:88-96
S10
Positions déclarées par le client
Mineur
Aucune vérification côté serveur.
LectureLes coordonnées reçues sont stockées et diffusées telles quelles.
engine/command.go:797-840
S11
Fermeture de canaux à l’arrêt
Mineur
shutdown() ferme les canaux des clients alors que les goroutines lecteur peuvent encore y écrire.
LectureRisque de panique à l’arrêt.
server/server.go:112-117 · server/client.go:86-97
Conformité au protocole
13 constats
ATTACK, TALK et QUESTS répondent dans un autre format que la RFC, et le README ne documente presque aucun écart.

P1
Format de ATTACK
Majeur
Attendu : OK {"attacker_hp":80,"target_hp":15,"damage":10,"status":"combat"}.
ReproduitReçu : OK {"player_action":{...},"npc_reactions":[...],"combat_state":"ONGOING"}.
engine/combat_system.go:53-58
P2
Format de TALK
Majeur
Attendu : OK Welcome to my bakery!.
ReproduitReçu : OK "Bienvenue dans Grow Harvest, aventurier !", chaîne JSON entre guillemets.
engine/command.go:624 · server/client.go:123-129
P3
Statuts de quête
Majeur
Attendu : active, completed.
ReproduitReçu : progress, unavailable. QUEST répond available alors que la quête est stockée en progress.
engine/command.go:440-451, 518
P4
Réponses hors ABNF
Majeur
Toute réponse commence par OK ou ERR <3 chiffres> <message>, tout événement par EVT <type> <données>.
ReproduitReçu : invalid coordinate format, invalid emote type, ERR Invitation already send, EVT new_leader=bob.
engine/command.go:810, 883 · engine/group.go:78, 211
P5
Codes d’erreur
Majeur
La grille demande les codes de la RFC §8.2 à l’identique.
Lecture14 codes ajoutés (200, 202, 400, 403, 407 à 414, 500, 902). 402, 901 et 406 servent chacun à deux erreurs différentes.
protocol/error_messages.go
P6
GROUP JOIN vers un groupe inexistant
Majeur
Une erreur métier (4xx) est attendue.
ReproduitReçu : ERR 500 INTERNAL_SERVER_ERROR.
engine/group.go:145-147
P7
STATUS sans état de combat
Majeur
Le sujet demande la santé et l’état de combat.
Reproduit"status":"normal" en plein combat ; le champ ne change jamais.
engine/command.go:336-347
P8
Écarts non documentés
Majeur
Le sujet exige que tout écart figure dans « Protocol Implementation ».
LectureAbsents du README : USE, FLEE, COMPLETE_QUEST, INSPECT, GROUPED, UNGROUPED, COMBAT STATS, CHAT COMBAT, les émotes, les événements ITEM, QUEST, COMBAT, les codes ajoutés, et P1 à P3.
README.md:65-78
P9
Clients non interchangeables
Majeur
Un client doit fonctionner avec le serveur d’un autre groupe.
LectureLe contrôleur envoie de lui-même INSPECT, GROUPED, UNGROUPED, COMBAT STATS, GET_PLAYER_POSITIONS, GET_ITEM_POSITIONS, inconnus de la RFC.
client/controller/events.go · responses.go
P10
Identifiants de salle incohérents
Mineur
LOOK doit donner des identifiants cohérents.
Reproduit"id":"room.place_du_village" mais les sorties et MOVE donnent place_du_village.
engine/command.go:174, 223
P11
CHAT GROUP hors groupe
Mineur
Attendu : ERR 401 NOT_IN_GROUP.
ReproduitReçu : OK, message perdu.
engine/command.go:226-253
P12
Pas d’écho à l’émetteur
Mineur
L’exemple du sujet renvoie EVT GLOBAL CHAT alice ... à alice.
ReproduitL’émetteur ne reçoit ni son message ni EVT STATS à sa connexion.
engine/utils.go:74-83
P13
Séparateurs stricts
Mineur
Les arguments sont séparés par des espaces.
ReproduitMOVE north (deux espaces) → ERR 301 NO_EXIT ; tabulation → ERR 400.
engine/handle_commands.go:30
Journalisation
4 constats
Le format est bon, mais il manque le nom du joueur sur les commandes et une partie des réponses.

L1
Commandes journalisées sans joueur
Majeur
Le sujet demande le nom du joueur et les paramètres.
Reproduit{"msg":"Client request","command":"CONNECT quinn"} : ni nom, ni identifiant, ni IP.
server/server.go:55
L2
Réponses non journalisées
Majeur
Toutes les réponses et tous les codes d’erreur doivent l’être.
ReproduitAucune ligne pour OK connected. À la lecture, idem pour ERR 200, 201, 202, 902 et 901 SERVER_FULL.
engine/handle_commands.go:38-45
L3
Volume de logs
Mineur
La journalisation ne doit pas dégrader les performances.
Reproduit13 012 lignes et 2,3 Mo en 12 s pour 6 joueurs en mouvement, écrites en synchrone dans la boucle du moteur.
engine/handle_commands.go:108-118
L4
Journal d’abus sans pseudo
Mineur
Le champ name reçoit la structure *Client.
Lectureslog.Warn("Command spam", "name", s.clients[cli.id], ...).
server/client.go:91
Monde et clients
4 constats
Le monde respecte les minimums. Les constats sur le GUI viennent de la lecture du code.

G1
GUI sans vue journal
Majeur
Le sujet demande une vue chat séparée de la vue journal.
LectureAppendServerResponse, AppendCliResponse, AppendCliMessage et UpdateDatas sont des fonctions vides.
client/gui/src/core/app_update.go:328-331
G2
Boutons manquants dans le GUI
Mineur
Boutons attendus pour LOOK, MOVE, TAKE, DROP, TALK, ATTACK, STATUS, QUEST, QUESTS, WHO, GROUP, QUIT.
LecturePas de bouton LOOK (envoi automatique) ; MOVE se fait en marchant vers un bord de carte. À confirmer à l’écran.
client/gui/src/updater/buttons.go
G3
État de groupe optimiste
Mineur
Le bouton bascule InGroup avant la réponse du serveur.
LectureEn cas d’erreur serveur, l’affichage diverge de l’état réel.
client/gui/src/updater/buttons.go:208
G4
Rôle de roger incohérent
Mineur
Les rôles de PNJ doivent être distincts.
Reproduitrole: dialogue mais porteur de la quête repare_barrage.
src/world.json
Architecture et qualité du code
12 constats
La séparation serveur, moteur, contrôleur et interfaces est claire. Le défaut de fond est le couplage par canaux bornés dans les deux sens.

A1
Cycle d’envois bloquants
Majeur
Le serveur pousse vers le moteur et le moteur vers le serveur, avec des envois bloquants des deux côtés. Un message diffusé à N joueurs crée N écritures.
ReproduitDémontré par B2.
server/server.go:52-66 · engine/utils.go
A2
Session et joueur désynchronisés
Majeur
sessions[id] survit à QUIT jusqu’à la fermeture de la socket, et aucune commande ne vérifie que le joueur existe encore.
ReproduitDémontré par B1.
engine/handle_commands.go:25-46
A3
Monde couplé à des constantes
Mineur
Salle de départ place_du_village, objet par défaut branche_solide, liste ValidMaps. Le validateur valid_room_type n’est jamais utilisé.
LectureConstantes en dur dans le moteur.
engine/valid_values.go:10-44 · protocol/protocol.go:4
A4
Le GUI dépend du paquet TUI
Mineur
Le contrôleur commun importe client/tui/panels.
LectureImport en tête de fichier.
client/controller/controller.go:9
A5
Réponses associées par une file
Mineur
Une ligne non sollicitée qui n’est pas un EVT décale toutes les réponses suivantes.
LectureFile lastCommands.
client/controller/controller.go:157-177
A6
Lecture client limitée à 64 Ko par ligne
Mineur
Une réponse plus longue coupe la connexion ; GET_ITEM_POSITIONS n’a pas de borne (S6).
Lecturebufio.Scanner par défaut.
client/network/client.go:103
A7
Code mort
Mineur
Canal logs, type Log, t_start, MaxFloodLimit, types Ressource et Consumable, removeItemFromWorld. Les effets mana, max_hp, status, initiative consomment l’objet sans rien faire.
LectureSymboles jamais référencés.
server/server.go:12-17 · engine/items.go:25-36 · engine/command.go:748-761
A8
Balise JSON erronée
Mineur
InCombat bool porte json:"omitempty" : le champ s’appelle donc omitempty.
LectureDéclaration de la structure Npc.
engine/npcs.go:16
A9
Nommage hétérogène
Mineur
Get_map, inform_room, end_combat à côté de handleCmdLook. Le sujet demande les conventions du langage.
LectureMélange de snake_case et de camelCase.
engine/*.go
A10
6 fichiers non formatés
Mineur
gofmt n’est pas appliqué partout.
Reproduitgofmt -l src liste 6 fichiers.
engine/combat_system.go · client/gui (5 fichiers)
A11
Tri d’initiative instable
Mineur
À initiative égale, l’ordre des joueurs varie d’un combat à l’autre.
Lecturesort.Slice au lieu de sort.SliceStable.
engine/combat_system.go:273-279
A12
Codes de sortie et option --ip
Mineur
Sortie avec le code 0 sur une erreur d’arguments ; --ip refuse les noms d’hôte.
Lecturereturn sans os.Exit(1) ; net.ParseIP seulement.
cmd/get_connection.go:66-70 · cmd/server/main.go:52-56
Build, Makefile et dépôt
9 constats
Toutes les cibles exigées existent et passent.

D1
Fichiers embarqués hors des dépendances
Majeur
Modifier src/world.json, une carte ou une image ne déclenche aucune recompilation. Piège direct pour l’exercice de modification en soutenance.
Reproduittouch src/world.json puis make exec/server → « is up to date ». Idem pour maps/jardin.json et exec/gui.
Makefile:36-41
D2
Cible par défaut et go mod tidy
Mineur
make seul ne lance que install, qui exécute go mod tidy (réseau requis, peut réécrire go.mod). build la relance à chaque fois.
Reproduitmake -n.
Makefile:44-47, 64
D3
Lint limité à go vet
Mineur
Ni gofmt, ni staticcheck : A10 passe inaperçu.
Reproduitmake lint passe malgré 6 fichiers non formatés.
Makefile:98-100
D4
Pas de .gitignore
Mineur
exec/ apparaît en fichiers non suivis après chaque build. La grille arrête l’évaluation en cas de fichier non autorisé.
Reproduitgit status après make build.
racine du dépôt
D5
Dépôt de 127 Mo
Mineur
Dont 117 Mo d’historique Git et 406 fichiers d’images.
Reproduitdu -sh.
src/client/gui/assets
D6
Packs graphiques sans crédit
Mineur
Dont un dossier « Sprout Lands - Sprites - premium pack ». Droit de redistribution à vérifier.
LectureAucune licence ni mention dans le README.
src/client/gui/assets
D7
Fichiers de travail commités
Mineur
prev_barrage.tmx, prev_foret_malicieuse.tmx, prev_riviere.tmx, sprout-lands.tiled-session.
LecturePrésents dans git ls-files.
src/client/gui/maps · tilesets
D8
Go 1.26.0 exigé
Mineur
Sur un poste plus ancien, Go télécharge une autre chaîne d’outils ou échoue hors ligne.
LectureDirective go 1.26.0.
go.mod:3
D9
Historique Git peu soigné
Mineur
Messages « Update » (5), « Add » (3), « Merged » (3), « Ficed bugs » ; 5 identités pour 2 auteurs ; 13 branches distantes non nettoyées.
Reproduitgit log, git shortlog.
historique
Tests
4 constats
La suite passe dans un terminal, mais échoue sans terminal et sous le détecteur de concurrence.

T1
make test échoue hors terminal
Majeur
TestTUINormalDisconnection ouvre /dev/tty puis panique. make check échoue donc en intégration continue.
ReproduitFAIL tap/src/tests/leaks (close of nil channel) ; ok dans un pseudo-terminal.
src/tests/leaks/normal_disconnection_test.go:28
T2
go test -race échoue
Majeur
Le Makefile ne lance jamais les tests avec -race.
ReproduitDATA RACE sur TestAntiFloodProtection (voir S5).
Makefile:84-89
T3
Zones non couvertes
Mineur
Aucun test pour B1, B2, B4 à B6, B8 ni le chargeur de monde. engine, server et protocol n’ont aucun test unitaire.
Lecture20 fonctions de test, toutes sous src/tests.
src/tests
T4
Monde de test différent du monde livré
Mineur
Le monde réel n’est jamais chargé par un test.
Reproduitdiff src/world.json src/tests/world.json : fichiers différents.
src/tests/world.json
README
9 constats
Plusieurs sections décrivent un comportement que le code n’a pas.

R1
« Group Contributions » vide
Critique
Le sujet demande les responsabilités de chaque membre.
ReproduitDeux lignes # TODO, rendues comme des titres de niveau 1.
README.md:178-182
R2
Système de quêtes décrit à tort
Majeur
États annoncés not_started, in_progress, completed ; récompense « automatique » ; objet gold_acorn.
ReproduitCode : available, progress, unavailable ; il faut COMPLETE_QUEST, jamais cité ; gold_acorn n’existe pas.
README.md:103-111
R3
Système de combat décrit à tort
Majeur
Réapparition annoncée à « Entrance (Overgrown Path) » ; formule « Attacker Damage - Target Shield ».
ReproduitCode : réapparition à place_du_village ; bouclier consommé côté joueur, ignoré côté PNJ ; USE absent ; aucun choix justifié.
README.md:82-99
R4
Journalisation décrite à tort
Majeur
Exemple avec timestamp, message, datas, client_id, COMMAND_EXECUTED.
ReproduitLe serveur écrit time, msg et des champs à plat. Le « Log Streamer » asynchrone du README serveur n’existe pas.
README.md:153-172 · src/server/README.md
R5
Choix d’interface du CLI non documenté
Majeur
Le sujet demande de dire si le CLI envoie la syntaxe RFC ou la traduit.
LectureAbsent du README racine.
README.md
R6
Pas de section « Instructions »
Mineur
Section de premier niveau attendue.
LectureSeulement un sous-titre dans « Building and Running ».
README.md:192
R7
Champ datas inexistant
Mineur
Le README parle d’un champ datas dans les réponses.
ReproduitSur le réseau, le JSON suit directement OK.
README.md:73-74
R8
Détails manquants ou faux
Mineur
Options --ip, --port et ARGS non documentées ; PNJ guide absent de « World Design » ; lien RFC sans URL ; serveur dit « persistent ».
LectureLe README du GUI cite aussi un bouton « Defend » inexistant.
README.md · src/client/gui/README.md:39
R9
Troisième auteur Git non mentionné
Mineur
Leo Goderneaux : 1 commit de fusion le 1er septembre 2026.
Reproduitgit shortlog -sne.
commit 5267e20
Recommandations, par ordre de priorité
Les quatre premières lèvent tous les points critiques. Elles sont à faire avant toute soutenance.

Corriger B1
Dans handleCommands, répondre ERR 200 NOT_CONNECTED quand e.players[pseudo] est nil, et vider e.sessions[id] dans playerQuits.
Casser le cycle B2 et B3
Le moteur ne doit jamais bloquer sur ServerOutput : file non bornée côté moteur, ou envoi non bloquant avec déconnexion du client lent. Fusionner JoinChan, LeaveChan et ServerInput en un seul canal ordonné.
Limiter les commandes de position (S4)
Leur donner un seau de jetons propre, par exemple 120 par seconde, au lieu d’une exemption totale.
Remplir « Group Contributions » (R1)
Et corriger les sections fausses du README : R2, R3, R4, R5, P8.
Valider les entrées (S1 à S3)
Nom limité à 32 caractères imprimables en UTF-8 valide ; rejeter par ERR 400 tout message contenant un caractère de contrôle.
Supprimer l’accès concurrent (S5)
Ne plus lire s.clients depuis la goroutine lecteur ; journaliser cli.id.
Aligner le protocole (P1 à P7)
Revenir aux formats de la RFC pour ATTACK, TALK, QUESTS et STATUS, ou les documenter un par un. Un code unique par erreur, toutes les réponses préfixées.
Réparer les groupes (B4 à B6)
Remettre promotion et send_promotion à zéro à l’exclusion et au départ, interdire l’auto-exclusion, n’utiliser que group.leader comme référence.
Réparer quêtes et combats (B8, B9, B12, B13)w
Filtrer sur progress dans refreshQuestProgress, cloner les objets donnés en récompense, refuser MOVE en combat, retirer le joueur de sa session à la déconnexion.
Compléter les logs (L1 à L3)
Journaliser commande et réponse au même endroit avec le nom du joueur ; passer les commandes de position au niveau DEBUG.
Fiabiliser l’outillage (D1, D3, D4, T1, T2)
Ajouter world.json, les cartes et les images aux dépendances du Makefile, gofmt -l à lint, un .gitignore, les tests avec -race, et ignorer le test TUI sans terminal.
Durcir le chargeur de monde (B18 à B20)
Sortir de la boucle sur toute erreur autre que io.EOF, remonter les clés dupliquées, valider les références des quêtes et la salle de départ.
Ce qui tient, et comment l’audit a été mené
Conforme et vérifié
Build : make install en 4,5 s, make build en 48 s, lint et clean sans erreur.
Message d’accueil, CONNECT, LOOK, MOVE, WHO, CHAT, TAKE, DROP, INVENTORY, QUEST, QUIT, groupes de base.
Objets pris par identifiant ou par nom affiché en plusieurs mots ; un objet pris disparaît de la salle.
Fragmentation et regroupement TCP, CRLF, Unicode valide, plafond de 50 connexions.
Combat : 100 PV au départ, riposte des PNJ, réapparition à 50 PV sur la place du village.
Monde : 12 salles en boucles, 3 impasses, 3 rôles de PNJ, 11 objets, 4 quêtes de 2 types.
Logs JSON avec niveaux, horodatage et IP.
Méthode
Références : rfc.html (protocole), tap.pdf (sujet v1.4), tap.html (grille de 17 rubriques). Aucun audit antérieur utilisé.
Exécuté sur une copie du dépôt : cibles du Makefile, go vet, gofmt -l, go test -race.
Serveur piloté par des scripts de sockets : toutes les commandes de la RFC, erreurs, groupes, quêtes, combats, charge.
Limite : les clients TUI et GUI n’ont pas été lancés en interactif. Leurs constats viennent de la lecture du code et sont marqués « Lecture ».
tap.html nomme une autre équipe (mdourdoi, mobenhab) ; seule sa grille est utilisée.
Critique : arrêt ou blocage du serveur, ou échec éliminatoire. Majeur : exigence du sujet ou de la RFC non tenue. Mineur : défaut sans effet bloquant.