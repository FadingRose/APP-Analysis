from lxml import etree
import sys
import json

def get_endpoints_for_app(endpointsPath, uniquePath):
    ends = []
    unique = []
    try:
        with open(endpointsPath,'r') as f:
            ends = f.readlines()
            # remove the '\n' from the end of each line
            ends = [end[:-1] for end in ends]
    except:
        print(f'Error in openning {endpointsPath}')
        exit(1)
    
    try:
        with open(uniquePath,'r') as f:
            unique = f.readlines()
            # remove the '\n' from the end of each line
            unique = [url[:-1] for url in unique]
    except:
        print(f'Error in openning {uniquePath}')
        exit(1)

    return ends, unique
    

def extract_apk_data_xml(manifestPath, endpointsPath, uniquePath):
    try:
        print(f'Parsing {manifestPath}')
        tree = etree.parse(manifestPath)
        root = tree.getroot()

        usr_permissions = root.findall('.//uses-permission')

        usr_permissions = [permission.attrib['{http://schemas.android.com/apk/res/android}name'] for permission in usr_permissions]

        endpoints, unique_urls = get_endpoints_for_app(endpointsPath,uniquePath)

        jsonData = {
                    'permissions': usr_permissions,
                    'package': root.attrib['package'],
                    'endpoints': endpoints,
                    'unique_urls': unique_urls,
                    'Activities': [activity.attrib['{http://schemas.android.com/apk/res/android}name'] for activity in root.findall('.//activity')],
                    'Services': [service.attrib['{http://schemas.android.com/apk/res/android}name'] for service in root.findall('.//service')],
                    'Receivers': [receiver.attrib['{http://schemas.android.com/apk/res/android}name'] for receiver in root.findall('.//receiver')],
                    'Providers': [provider.attrib['{http://schemas.android.com/apk/res/android}name'] for provider in root.findall('.//provider')],
        }

        with open('target.json','w') as f:
            json.dump(jsonData,f)
    except Exception as e:
        print(f'Error in parsing {manifestPath}')
        print(e)
        exit(1)

if __name__ == '__main__':
    # receive cmd line arguments
    manifest = sys.argv[1]
    endpoints = sys.argv[2]
    unique = sys.argv[3]
    extract_apk_data_xml(manifest,endpoints,unique)