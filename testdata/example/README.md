# Beispielprojekt

Dieses Projekt ist ein Beispielprojekt und zeigt die Verwendung verschiedener Annotationen.


## Zeiterfassungsmanagement

Package domain enthält den Bounded Context über die [Zeiterfassung](#zeiterfassung).

### Anwendungsfälle

#### Aufstehen

hello


#### Aufstehen in der Zeiterfassung

[Aufstehen](#aufstehen) [Zeiterfassung](#zeiterfassung)smethode


#### Aufstehen woanders

[Aufstehen](#aufstehen) [Zeiterfassung](#zeiterfassung)smethode


#### Beschwerde einreichen

Dieser Anwendungsfall ist noch nicht dokumentiert.

#### Zeiten loggen

Cooles Zeitbuchen ist angesagt.


### Werte

#### Zeitlog

Dieser Werttyp ist noch nicht dokumentiert.

### Entitäten

#### Mitarbeiter

[Mitarbeiter](#mitarbeiter) arbeitet bei seinem Arbeitgeber.


### Aggregate

#### User

Dieses Aggregat ist noch nicht dokumentiert.

### Domänenservices

#### Zeiterfassung

Dieser Service ist noch nicht dokumentiert.

### Repositories

#### Zeitaufzeichnungen

[Zeitlog](#zeitlog)Repo manages the [Zeitlog](#zeitlog)s.


## Berechtigungskonzept

Im Folgenden werden alle auditierten Berechtigungen dargestellt.
Diese Berechtigungen sind Aktor-gebunden, d.h. ein Nutzer oder Drittsysteme müssen diese Rechte zugewiesen bekommen haben, um den Anwendungsfall ausführen zu dürfen.
Achtung: es kann dynamische bzw. objektbezogene Rechte in Anwendungsfällen geben, die unabhängig von Berechtigungen das Darstellen, Löschen oder Bearbeiten von vertraulichen Informationen erlaubt. Diese sind hier nicht erfasst, sondern sind in der jeweiligen Dokumentation der Anwendungsfälle erwähnt.

|Berechtigung|Anwendungsfall|
|----|----|
|de.worldiety.aufstehen|[Aufstehen](#aufstehen)|
|de.worldiety.aufstehen2|[Zeiten loggen](#zeiten-loggen)|
|de.worldiety.woanders.aufstehen|[Aufstehen woanders](#aufstehen-woanders)|
|de.worldiety.zeiterfassung.aufstehen|[Aufstehen in der Zeiterfassung](#aufstehen-in-der-zeiterfassung)|

Die folgenden Anwendungsfälle sind grundsätzlich ohne Autorisierung verwendbar, erfordern also keine Berechtigungen und werden auch nicht auditiert.

|Berechtigung|Anwendungsfall|
|----|----|
|jeder|[Beschwerde einreichen](#beschwerde-einreichen)|
