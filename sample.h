#ifdef INTERNAL
typedef struct Sample Sample;
#else
typedef void Sample;
#endif
typedef int (*SampleCallback) (Sample* sample, int number, void* arbitrary_data);
Sample* create_sample(int number);
int destroy_sample(Sample* sample);
int invoke_callback(Sample* sample);
int register_callback(Sample* sample, SampleCallback callback, void* arbitrary_data);
