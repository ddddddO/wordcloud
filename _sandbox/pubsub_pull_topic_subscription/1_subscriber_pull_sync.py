from google.cloud import pubsub
from io import BytesIO
from PIL import Image

project_id = "test-emu-555"
topic_name = "my_topic"
subscription_name = "my_subscription"

subscriber = pubsub.SubscriberClient()
topic_path = "projects/{project}/topics/{topic}".format(project=project_id, topic=topic_name)
subscription_path = subscriber.subscription_path(project_id, subscription_name)
# subscriber.create_subscription(request={"name": subscription_path, "topic": topic_path})
print("block..")

response = subscriber.pull(
    request={
        "subscription": subscription_path,
        "max_messages": 1,
    }
)

print("block....")

for msg in response.received_messages:
    print("Received message:", msg.message.data)
    print("Received message:", msg.message)

    # img = Image.open(BytesIO(msg.message.data))
    # img.save('published_wordcloud.png')

print("pass")

ack_ids = [msg.ack_id for msg in response.received_messages]
subscriber.acknowledge(
    request={
        "subscription": subscription_path,
        "ack_ids": ack_ids,
    }
)