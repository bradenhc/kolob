# Kolob - Simple & Secure Accountless Collaboration

Kolob is a simple collaboration tool designed to target an audience that may not
have email addresses or mobile phones but does have access to the internet. Once
a verified user has created a group, other users can use the group information
to sign in to Kolob and start posting messages to the group.

## Motivation

TODO

## Capabilities

TODO

## Command Interface

The `kolob` command is used to launch a single Kolob server.

The `kolobctl` command is used to manage several kolob servers. It provides a
clean user interfaces that lets users create new groups and monitors the Kolob
server associated with a group. `kolobctl` uses containerization technologies to
do most of the heavy lifting, so you should make sure Docker is installed if you
are going to be using it.

## Data Model

Kolob focuses on a minimal feature set in order to provide the highest quality
experience for a special niche of users. The remainder of this section breaks
down this data model and explains the motivation behind the different elements.

### Groups

The central element in this data model is the **Group**. A group is where users
can join together and post messages about various topics. Every group has at
least one administrator called the Group Creator. Additional details about the
Group Creator are provided in following sections.

A single Kolob server can only run one group. This enables complete isolation
of group data and also make the project easier to maintain. As such, the idea of
a group is more conceptual than it is concrete in the program's implementation.

### Conversations

A **Conversation** is a time-ordered list of messages sent by group members
surrounding a particular topic. A single group may contain more than one
conversation. When a group is first created it contains a single conversation
titled "General" that serves as a starting point for the group. The Group
Creator is free to remove this conversation after the group is created so long
as there is at least one additional conversation in the group.

### Members

A **Member** belongs to one and only one group. Member's are identified within a
group by their username. A username is unique within a group, but Kolob does not
require that usernames be unique across groups.

#### Group Creators

The member that creates the group is called the **Group Creator**. While other
group members do not need to provide a separate email or username, the Group
Creator _must_ provide a phone number or email address and respond to a
confirmation before the group is created.

This security feature protects the Kolob server from being overwhelemed with
fake groups and helps provide group members with a sense of security because
they must know the group creator personally in order to join a group, as group
information must be given to them by a group creator.

### Messages

TODO

## Data Storage

Kolob can be extended to support multiple backend data storage technologies. The
default backend is driven by SQLite.

Data is always written to disk before it is applied to the in-memory store.
Data on disk is always encrypted.

## Security

All group information is encrypted using AES with a 256-bit key generated by the
PBKDF2 algorithm with a user-provided password, a 32 byte salt, and 1,000,000
iterations. The user provided password must be at least 16 characters and
contain at least one lowercase letter, one uppercase leter, one number, and one
special character.

> NOTE: The iteration count was selected based on the [OWASP suggestion] of
> 600,000 or more as referenced in a document of [comments on SP 800-132]
> provided to the NIST. The password criteria was selected based on the
> [password guidelines] provided by OWASP.

## Design

TODO

<!-- LINKS -->
<!-- prettier-ignore-start -->
[OWASP suggestion]: https://cheatsheetseries.owasp.org/cheatsheets/Password_Storage_Cheat_Sheet.html
[Comments on SP 800-132]: https://csrc.nist.gov/csrc/media/Projects/crypto-publication-review-project/documents/initial-comments/sp800-132-initial-public-comments-2023.pdf
[password guidelines]: https://cheatsheetseries.owasp.org/cheatsheets/Authentication_Cheat_Sheet.html
<!-- prettier-ignore-end -->
