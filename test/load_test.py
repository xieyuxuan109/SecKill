import requests
import threading
import time
import random

# 配置
BASE_URL = "http://localhost:8888"
ACTIVITY_ID = 1
THREAD_COUNT = 50
REQUESTS_PER_THREAD = 100
results = {"success": 0, "failed": 0}
lock = threading.Lock()

def seckill_request():
    global results
    for _ in range(REQUESTS_PER_THREAD):
        try:
            response = requests.post(
                f"{BASE_URL}/api/seckill",
                json={"activityId": ACTIVITY_ID},
                timeout=5
            )
            with lock:
                if response.status_code == 200:
                    results["success"] += 1
                else:
                    results["failed"] += 1
        except Exception as e:
            with lock:
                results["failed"] += 1

def create_activity():
    try:
        now = int(time.time())
        response = requests.post(
            f"{BASE_URL}/api/act/create",
            json={
                "name": "压测活动",
                "stock": 1000,
                "startAt": now - 300,
                "endAt": now + 300
            },
            timeout=10
        )
        print(f"创建活动响应: {response.status_code}")
        print(f"响应内容: {response.text}")
    except Exception as e:
        print(f"创建活动失败: {e}")

if __name__ == "__main__":
    print("=== 秒杀服务压测 ===")
    
    # 先创建活动
    print("1. 创建秒杀活动...")
    create_activity()
    time.sleep(2)
    
    # 开始压测
    print(f"2. 开始压测，线程数: {THREAD_COUNT}，每线程请求数: {REQUESTS_PER_THREAD}")
    threads = []
    start_time = time.time()
    
    for _ in range(THREAD_COUNT):
        t = threading.Thread(target=seckill_request)
        threads.append(t)
        t.start()
    
    for t in threads:
        t.join()
    
    end_time = time.time()
    total_requests = THREAD_COUNT * REQUESTS_PER_THREAD
    duration = end_time - start_time
    qps = total_requests / duration
    
    print("\n=== 压测结果 ===")
    print(f"总请求数: {total_requests}")
    print(f"成功数: {results['success']}")
    print(f"失败数: {results['failed']}")
    print(f"成功率: {results['success'] / total_requests * 100:.2f}%")
    print(f"总耗时: {duration:.2f}秒")
    print(f"QPS: {qps:.2f}")
