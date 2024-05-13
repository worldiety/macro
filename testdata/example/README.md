# example

## Zeiterfassungsmanagement

### Anwendungsfälle

#### <a name="Aufstehen"></a> Aufstehen

Dieser Anwendungsfall ist noch nicht dokumentiert.

#### <a name="BeschwerdeEinreichen"></a> Beschwerde einreichen

Dieser Anwendungsfall ist noch nicht dokumentiert.

#### <a name="ZeitBuchen"></a> ZeitBuchen

Cooles Zeitbuchen ist angesagt.


### Werte

#### Zeitlog

Dieser Werttyp ist noch nicht dokumentiert.

### Entitäten

#### Mitarbeiter

Mitarbeiter arbeitet bei seinem Arbeitgeber.


### Aggregate

#### User

Dieses Aggregat ist noch nicht dokumentiert.

### Domänenereignisse

### Domänenservices

#### Zeiterfassung

Dieser Service ist noch nicht dokumentiert.

## Berechtigungskonzept

Im Folgenden werden alle auditierten Berechtigungen dargestellt.
Diese Berechtigungen sind Aktor-gebunden, d.h. ein Nutzer oder Drittsysteme müssen diese Rechte zugewiesen bekommen haben, um den Anwendungsfall ausführen zu dürfen.
Achtung: es kann dynamische bzw. objektbezogene Rechte in Anwendungsfällen geben, die unabhängig von Berechtigungen das Darstellen, Löschen oder Bearbeiten von vertraulichen Informationen erlaubt. Diese sind hier nicht erfasst, sondern sind in der jeweiligen Dokumentation der Anwendungsfälle erwähnt.

|Berechtigung|Anwendungsfall|
|----|----|
|de.worldiety.aufstehen|[Aufstehen](#Aufstehen)|
|de.worldiety.aufstehen2|[ZeitBuchen](#ZeitBuchen)|
Die folgenden Anwendungsfälle sind grundsätzlich ohne Autorisierung verwendbar, erfordern also keine Berechtigungen und werden auch nicht auditiert.

|Berechtigung|Anwendungsfall|
|----|----|
|jeder|[Beschwerde einreichen](#BeschwerdeEinreichen)|
