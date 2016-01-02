from warehouse.models import Prizes

def getallprizes():
    return Prizes.objects.all().count()
