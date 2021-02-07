from google.cloud import pubsub_v1

project_id = "test-emu-555"
topic_name = "my_topic"
publisher = pubsub_v1.PublisherClient()
topic_path = publisher.topic_path(project_id, topic_name)
topic = publisher.create_topic(request={"name": topic_path})
print('Topic created: {}'.format(topic))