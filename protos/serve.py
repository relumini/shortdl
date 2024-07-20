import grpc
import aiofiles
import os
import httpx
import zipfile
from concurrent import futures
from lib.dtiktok.crawlers.hybrid.hybrid_crawler import HybridCrawler
import multi_pb2
import multi_pb2_grpc
import asyncio
from grpc import aio

HybridCrawler = HybridCrawler()


class DownloadShortServicer(multi_pb2_grpc.DownloadShortServicer):
    async def DownTiktok(self, request, context):
        try:
            data = await HybridCrawler.hybrid_parsing_single_video(url=request.url, minimal=True)
        except Exception as e:
            print(e)
            return multi_pb2.ReturnsReply(status="Failed")
        
        try:
            data_type = data.get('type')
            platform = data.get('platform')
            aweme_id = data.get('aweme_id')
            download_path = "../download"

            os.makedirs(download_path, exist_ok=True)

            if data_type == 'video' or data_type == None:
                file_name = f"{platform}_{aweme_id}.mp4"
                url = data.get('video_data').get('nwm_video_url_HQ')
                file_path = os.path.join(download_path, file_name)
                if os.path.exists(file_path):
                    return multi_pb2.ReturnsReply(status="already downloaded")
                __headers = await HybridCrawler.TikTokWebCrawler.get_tiktok_headers() if platform == 'tiktok' else await HybridCrawler.DouyinWebCrawler.get_douyin_headers()
                response = await self.fetch_data(url, headers=__headers)
                
                async with aiofiles.open(file_path, 'wb') as out_file:
                    await out_file.write(response.content)

                return multi_pb2.ReturnsReply(status="downloaded")

            elif data_type == 'image':
                zip_file_name = f"{platform}_{aweme_id}_images.zip"
                zip_file_path = os.path.join(download_path, zip_file_name)

                if os.path.exists(zip_file_path):
                    return multi_pb2.ReturnsReply(status="zip exist")

                urls = data.get('image_data').get('no_watermark_image_list')
                image_file_list = []
                for url in urls:
                    response = await self.fetch_data(url)
                    index = int(urls.index(url))
                    content_type = response.headers.get('content-type')
                    file_format = content_type.split('/')[1]
                    file_name = f"{platform}_{aweme_id}_{index + 1}.{file_format}"
                    file_path = os.path.join(download_path, file_name)
                    image_file_list.append(file_path)

                    async with aiofiles.open(file_path, 'wb') as out_file:
                        await out_file.write(response.content)

                with zipfile.ZipFile(zip_file_path, 'w') as zip_file:
                    for image_file in image_file_list:
                        zip_file.write(image_file, os.path.basename(image_file))

                return multi_pb2.ReturnsReply(status="downloaded")

        except Exception as e:
            print(e)
            return multi_pb2.ReturnsReply(status="error occurred")

    def DownYoutube(self, request, context):
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    async def fetch_data(self, url: str, headers: dict = None):
        headers = {
            'User-Agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36'
        } if headers is None else headers.get('headers')
        async with httpx.AsyncClient() as client:
            response = await client.get(url, headers=headers)
            response.raise_for_status()
            return response

async def serve():
    server = aio.server()
    multi_pb2_grpc.add_DownloadShortServicer_to_server(DownloadShortServicer(), server)
    listen_addr = 'localhost:50051'
    server.add_insecure_port(listen_addr)
    await server.start()
    print(f"Server started on {listen_addr}")
    await server.wait_for_termination()
if __name__ == '__main__':
    asyncio.run(serve())