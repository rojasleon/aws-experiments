import { SendMessageCommand } from '@aws-sdk/client-sqs';
import { Upload } from '@aws-sdk/lib-storage';
import { APIGatewayProxyHandler } from 'aws-lambda';
import { GetProductsEvent, ProductRequest } from '../../../common/types';
import { s3, sqs } from '../clients';

export const handler: APIGatewayProxyHandler = async (event) => {
  try {
    const { productId, reason }: ProductRequest = JSON.parse(event.body ?? '');

    if (!productId || !reason) {
      throw new Error('Must provide productId and message properties');
    }
    if (reason !== 'purchase' && reason !== 'question') {
      throw new Error('reason must be "purchase" or "question"');
    }

    const bucketPath = `requests/${(Math.random() + 1)
      .toString(36)
      .substring(7)}.json`;

    const multiPartUpload = new Upload({
      client: s3,
      params: {
        Bucket: process.env.BUCKET_NAME,
        Key: bucketPath,
        Body: JSON.stringify({ productId, reason }),
        ContentType: 'application/json'
      },
      queueSize: 4,
      partSize: 1024 * 1024 * 5,
      leavePartsOnError: false
    });

    await multiPartUpload.done();

    const eventMessage: GetProductsEvent = {
      bucketPath
    };

    await sqs.send(
      new SendMessageCommand({
        QueueUrl: process.env.REQUEST_CREATION_QUEUE_URL,
        MessageBody: JSON.stringify(eventMessage)
      })
    );
    return {
      statusCode: 200,
      body: JSON.stringify({ msg: 'hello there' })
    };
  } catch (error) {
    return {
      statusCode: 500,
      body: JSON.stringify({ error })
    };
  }
};
