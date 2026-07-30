"""
Fine-tuning model kecil (sentence-transformers) untuk klasifikasi intent
pertanyaan Customer Support: retur, pengiriman, pembayaran, dll.
"""
from sentence_transformers import SentenceTransformer, losses, InputExample
from torch.utils.data import DataLoader

MODEL_NAME = "all-MiniLM-L6-v2"
OUTPUT_DIR = "./finetuned-support-intent-model"


def build_training_data() -> list[InputExample]:
    return [
        InputExample(texts=["Kapan pesanan saya sampai?", "status pengiriman"], label=0.9),
        InputExample(texts=["Bagaimana cara retur barang?", "kebijakan retur"], label=0.95),
        InputExample(texts=["Kenapa pembayaran saya gagal?", "masalah pembayaran"], label=0.9),
        InputExample(texts=["Apakah bisa bayar COD?", "metode pembayaran"], label=0.85),
        InputExample(texts=["Barang saya rusak, bisa tukar?", "kebijakan retur"], label=0.9),
    ]


def train():
    model = SentenceTransformer(MODEL_NAME)
    train_examples = build_training_data()
    train_dataloader = DataLoader(train_examples, shuffle=True, batch_size=8)
    train_loss = losses.CosineSimilarityLoss(model)

    model.fit(
        train_objectives=[(train_dataloader, train_loss)],
        epochs=3,
        warmup_steps=10,
        output_path=OUTPUT_DIR,
    )
    print(f"[finetune] Model tersimpan di {OUTPUT_DIR}")


if __name__ == "__main__":
    train()