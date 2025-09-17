import grpc
from concurrent import futures
import time
import torch
from transformers import AutoTokenizer, AutoModelForSequenceClassification
import nlp_pb2
import nlp_pb2_grpc

class NLPServiceServicer(nlp_pb2_grpc.NLPServiceServicer):
    def __init__(self):
        MODEL_NAME = "unitary/toxic-bert"
        self.tokenizer = AutoTokenizer.from_pretrained(MODEL_NAME)
        self.model = AutoModelForSequenceClassification.from_pretrained(MODEL_NAME)

    def Analyze(self, request, context):
        text = request.text
        inputs = self.tokenizer(text, return_tensors="pt", truncation=True)
        with torch.no_grad():
            logits = self.model(**inputs).logits
            probs = torch.sigmoid(logits)[0]
        resp = nlp_pb2.AnalyzeResponse(
            toxic=bool(probs[0] > 0.5),
            profanity=bool(probs[2] > 0.5),
            violence=bool(probs[3] > 0.5),
            racial=bool(probs[5] > 0.5),
            score=float(probs[0])
        )
        return resp

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=4))
    nlp_pb2_grpc.add_NLPServiceServicer_to_server(NLPServiceServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    print("NLP gRPC server running on port 50051")
    try:
        while True:
            time.sleep(86400)
    except KeyboardInterrupt:
        server.stop(0)

if __name__ == "__main__":
    serve()
