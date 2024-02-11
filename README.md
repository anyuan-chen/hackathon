Tables:

user:
name
company
email
phone
role
hashed_secret

skills:
user_id
skill
rating

user_history:
date_change
who_changed_id
id (fk)
name
company
email
phone
role
hashed_secret

skills_changes:
date_change
who_changed_id
user_id
skill_rating

all users:
(accessible to only admin keys)

get/update user info:
(accessible to users themselves and admins)

deploy to fly.io
