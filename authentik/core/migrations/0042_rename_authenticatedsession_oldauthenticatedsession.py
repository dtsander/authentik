# Generated by Django 5.0.10 on 2024-12-18 14:18

from django.db import migrations


class Migration(migrations.Migration):

    dependencies = [
        ("authentik_core", "0041_applicationentitlement"),
        ("authentik_providers_oauth2", "0026_alter_accesstoken_session_and_more"),
        ("authentik_providers_rac", "0005_alter_racpropertymapping_options"),
    ]

    operations = [
        migrations.RenameModel(
            old_name="AuthenticatedSession",
            new_name="OldAuthenticatedSession",
        ),
    ]
